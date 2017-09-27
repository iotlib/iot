package httpserver

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/twinone/iot/backend/model"
)

const (
	userInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	defaultCookie    = "default"
)

func (s *Server) dashboardHandler(w http.ResponseWriter, r *http.Request, cookie *sessions.Session, user *model.User) {
	http.ServeFile(w, r, "www/index.html")
}
