package httpserver

import (
	"net/http"
	"html/template"

	"github.com/namsral/flag"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"encoding/base64"
	"crypto/rand"
	"github.com/gorilla/sessions"
	"context"
	"log"
	"io/ioutil"
	"encoding/json"
)

const (
	userInfoEndpoint       = "https://www.googleapis.com/oauth2/v3/userinfo"
	defaultCookieStoreName = "default"
)

var (
	oAuthCallbackUrl  = flag.String("callback_url", "", "OAuth Callback URL")
	oAuthClientID     = flag.String("client_id", "", "OAuth Client ID")
	oAuthClientSecret = flag.String("client_secret", "", "OAuth Client Secret")
	cookieStoreSecret = flag.String("cookie_store_secret", "", "Cookie-store secret")
)
var homeTemplate = template.Must(template.ParseFiles("www/index.tmpl"))

var store *sessions.CookieStore
var cfg *oauth2.Config

type User struct {
	Sub string `json:"sub"`
	Name string `json:"name"`
	GivenName string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Profile string `json:"profile"`
	Picture string `json:"picture"`
	Email string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Gender string `json:"gender"`
}



func Setup() {
	store = sessions.NewCookieStore([]byte(*cookieStoreSecret))
	cfg = &oauth2.Config{
		ClientID:     *oAuthClientID,
		ClientSecret: *oAuthClientSecret,
		RedirectURL:  *oAuthCallbackUrl,
		Scopes: []string{
			// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func getLoginURL(state string) string {
	return cfg.AuthCodeURL(state)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, defaultCookieStoreName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	tok, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Got access token:", tok)
	client := cfg.Client(context.Background(), tok)
	resp, err := client.Get(userInfoEndpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var u = &User{}
	json.Unmarshal(data, u)

	w.Write([]byte("Logged in as: " + u.Email))
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, defaultCookieStoreName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	state := randToken()
	session.Values["state"] = state

	session.Save(r, w)
	homeTemplate.Execute(w, getLoginURL(state))
}
