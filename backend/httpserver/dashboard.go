package httpserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/twinone/iot/backend/model"
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

func (s *Server) dashboardHandler(w http.ResponseWriter, r *http.Request, cookie *sessions.Session, user *model.User) {
	di := &DashboardInfo{
		User:    user,
		Devices: s.hub.GetDevices(user.Email),
	}
	log.Println("devices:", di.Devices)

	// Logged in
	dashboardTemplate.Execute(w, di)
}
