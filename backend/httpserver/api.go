package httpserver

import (
	"net/http"

	"encoding/json"
	"log"

	"github.com/gorilla/sessions"
	"github.com/twinone/iot/backend/model"
)

func (s *Server) apiHandler(w http.ResponseWriter, r *http.Request, cookie *sessions.Session, user *model.User) {
	cmd := r.URL.Query().Get("cmd")
	if cmd != "" {
		for _, conn := range s.hub.GetConns(user.Email) {
			log.Println("Sending cmd:", cmd)
			conn.Send <- []byte(cmd)
		}
		return
	}
	di := &model.DashboardInfo{
		User:    user,
		Devices: s.hub.GetDevices(user.Email),
	}
	WriteJSON(w, di)
}

func WriteJSON(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Println("Error marshaling json:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
