package httpserver

import (
	"net/http"
	"html/template"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"encoding/base64"
	"crypto/rand"
	"github.com/gorilla/sessions"
	"context"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/twinone/iot/backend/model"
	"github.com/twinone/iot/backend/ws"
)

const (
	userInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	defaultCookie    = "default"
)

var signinTemplate = template.Must(template.ParseFiles("www/signin.tmpl"))
var dashboardTemplate = template.Must(template.ParseFiles("www/dashboard.tmpl"))

type DashboardInfo struct {
	User    *model.User
	Devices []*model.Device
}

type Server struct {
	store *sessions.CookieStore
	cfg   *oauth2.Config
	hub   *ws.Hub
}

func New(config map[string]*string, hub *ws.Hub) (s *Server) {
	log.Println("Config:", config)
	return &Server{
		hub:   hub,
		store: sessions.NewCookieStore([]byte(*config["cookie_store_secret"])),

		cfg: &oauth2.Config{
			ClientID:     *config["client_id"],
			RedirectURL:  *config["callback_url"],
			ClientSecret: *config["client_secret"],
			Scopes: []string{
				// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (s *Server) getLoginURL(state string) string {
	return s.cfg.AuthCodeURL(state)
}

func (s *Server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.GetSession(w, r)
	if err != nil {
		return
	}
	savedState := session.Values["state"]
	state := r.URL.Query()["state"][0]
	if savedState != state {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	code := r.URL.Query()["code"][0]
	log.Println("Got auth code:", code)

	tok, err := s.cfg.Exchange(context.Background(), code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Got access token:", tok)
	client := s.cfg.Client(context.Background(), tok)
	resp, err := client.Get(userInfoEndpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var u = &model.User{}
	json.Unmarshal(data, u)

	session.Values["profile"] = data
	log.Println("profile:", string(data))
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Server) IndexLoginDashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, user, err := s.GetUser(w, r)
	if err != nil {
		return
	}
	if user == nil {
		state := randToken()
		session.Values["state"] = state

		session.Save(r, w)
		signinTemplate.Execute(w, s.getLoginURL(state))
		return
	}

	di := &DashboardInfo{
		User:    user,
		Devices: s.hub.GetDevices(user.Email),
	}
	log.Println("devices:", di.Devices)

	// Logged in
	dashboardTemplate.Execute(w, di)
}

func (s *Server) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/", s.IndexLoginDashboardHandler)
	r.HandleFunc("/signout", s.SignOutHandler)
	r.HandleFunc("/auth/callback", s.AuthHandler)
}

func (s *Server) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	session, _, err := s.GetUser(w, r)
	if err != nil {
		return
	}
	session.Values["state"] = ""
	session.Values["profile"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Server) GetSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	session, err := s.store.Get(r, defaultCookie)
	if err != nil {
		w.Header().Add("Set-Cookie", defaultCookie+"=empty;Max-Age=0")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error, try refreshing the page"))
		log.Println("Error getting session:", err.Error())
		return nil, err
	}
	return session, nil
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) (*sessions.Session, *model.User, error) {
	session, err := s.GetSession(w, r)
	if err != nil {
		return nil, nil, err
	}

	profile, ok := session.Values["profile"]
	if !ok {
		return session, nil, nil
	}
	data, ok := profile.([]byte)
	if !ok {
		return session, nil, nil
	}
	var u = &model.User{}
	err = json.Unmarshal(data, u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting user:", err.Error())
		return session, nil, err
	}
	return session, u, nil
}
