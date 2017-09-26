package httpserver

import (
	"encoding/json"
	"net/http"

	"context"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"log"

	"github.com/gorilla/sessions"
	"github.com/twinone/iot/backend/db"
	"github.com/twinone/iot/backend/model"
)

type AuthedHandler = func(w http.ResponseWriter, r *http.Request, c *sessions.Session, user *model.User)

func (s *Server) Auth(next AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := s.GetCookie(r)
		u := s.GetUser(c)
		if u == nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		next(w, r, c, u)
	}
}

func (s *Server) signinHandler(w http.ResponseWriter, r *http.Request) {
	c := s.GetCookie(r)
	uid := randToken()
	if id, ok := c.Values["state"]; ok {
		uid = id.(string)
	}
	c.Values["state"] = uid
	c.Save(r, w)

	log.Println("uid:", uid)
	signinTemplate.Execute(w, s.getLoginURL(uid))
}

func (s *Server) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	cookie := s.GetCookie(r)
	savedState, ok := cookie.Values["state"]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	state := r.URL.Query().Get("state")
	if savedState != state {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	tok, err := s.cfg.Exchange(context.Background(), code)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//log.Println("Got access token:", tok)
	client := s.cfg.Client(context.Background(), tok)
	resp, err := client.Get(userInfoEndpoint)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var authedUser = &model.User{}
	json.Unmarshal(data, authedUser)

	dbUser := db.GetUserByEmail(authedUser.Email)
	if dbUser == nil {
		db.InsertUser(authedUser)
	}

	db.InsertAccessToken(&model.AccessToken{
		Email: authedUser.Email,
		Token: state,
	})

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

func (s *Server) signOutHandler(w http.ResponseWriter, r *http.Request, cookie *sessions.Session, user *model.User) {
	tok := cookie.Values["state"]
	log.Println("Removing tok:", tok)
	db.RemoveAccessToken(tok.(string))

	// Generate a new state for the user
	cookie.Values["state"] = randToken()
	cookie.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Returns a User if the request is authenticated or nil if not
func (s *Server) GetUser(cookie *sessions.Session) *model.User {
	tok, ok := cookie.Values["state"]
	if !ok {
		return nil
	}
	return db.GetUserByAccessToken(tok.(string))
}

func (s *Server) getLoginURL(state string) string {
	return s.cfg.AuthCodeURL(state)
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
