package httpserver

import (
	"net/http"

	"github.com/Don-V/mongostore"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/twinone/iot/backend/db"
	"github.com/twinone/iot/backend/ws"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Server struct {
	store *mongostore.MongoStore
	cfg   *oauth2.Config
	hub   *ws.Hub
}

func New(config map[string]*string, hub *ws.Hub) (s *Server) {
	return &Server{
		hub: hub,
		//store: sessions.NewCookieStore([]byte(*config["cookie_store_secret"])),
		store: mongostore.NewMongoStore(
			db.GetCookieCollection(),
			3600*24*365*10,
			false,
			[]byte(*config["cookie_store_secret"])),

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

func (s *Server) GetCookie(r *http.Request) *sessions.Session {
	session, _ := s.store.Get(r, defaultCookie)
	return session
}

func (s *Server) RegisterHandlers(r *mux.Router) {

	//	r.HandleFunc("/", s.indexHandler)
	r.HandleFunc("/signin", s.signinHandler)
	r.HandleFunc("/auth/callback", s.authCallbackHandler)

	// protected endpoints
	apiRouter := r.PathPrefix("/api/").Subrouter()
	s.registerApiHandlers(apiRouter)

	r.HandleFunc("/signout", s.Auth(s.signOutHandler))
	r.HandleFunc("/dashboard", s.Auth(s.dashboardHandler))

	r.HandleFunc("/", s.indexHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/")))

}
