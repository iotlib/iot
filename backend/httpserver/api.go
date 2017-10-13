package httpserver

import (
	"net/http"

	"encoding/json"
	"log"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/twinone/iot/backend/db"
	"github.com/twinone/iot/backend/model"
)

func (s *Server) execHandler(w http.ResponseWriter, r *http.Request, cookie *sessions.Session, user *model.User) {
	decoder := json.NewDecoder(r.Body)
	var e model.Execution
	err := decoder.Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Hey")

	for _, conn := range s.hub.GetConns(user.Email) {
		if conn.Device.Id == e.DeviceId {
			log.Println("Sending cmd", e.Cmd, "to", e.DeviceId)
			conn.Send <- []byte(e.Cmd)
		}
	}
	return

}

func (s *Server) profileHandler(w http.ResponseWriter, r *http.Request, cookie *sessions.Session, user *model.User) {
	di := &model.DashboardInfo{
		User:      user,
		Devices:   s.hub.GetDevices(user.Email),
		Functions: db.FindFunctionsByEmail(user.Email),
	}
	WriteJSON(w, di)
}

func (s *Server) functionHandler(w http.ResponseWriter, r *http.Request, cookie *sessions.Session, user *model.User) {
	decoder := json.NewDecoder(r.Body)
	var f model.Function
	err := decoder.Decode(&f)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Decoded:", f)
	defer r.Body.Close()
	if len(f.Cmd) > 20 ||
		f.Cmd == "" ||
		f.Pin < 0 ||
		f.Pin > 30 ||
		len(f.Name) > 32 ||
		f.Name == "" {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f.Owner = user.Email
	id := db.InsertFunction(&f)

	WriteJSON(w, map[string]string{
		"id": id,
	})
}

func (s *Server) deleteFunctionHandler(w http.ResponseWriter, r *http.Request, cookie *sessions.Session, user *model.User) {
	id := mux.Vars(r)["id"]
	email := user.Email

	log.Println("deleted function:", id, email)
	db.RemoveFunction(id, email)

}

func (s *Server) registerApiHandlers(r *mux.Router) {
	r.Handle("/profile", s.Auth(s.profileHandler)).Methods("GET")
	r.Handle("/function", s.Auth(s.functionHandler)).Methods("POST", "GET")
	r.Handle("/function/{id}", s.Auth(s.deleteFunctionHandler)).Methods("DELETE")
	r.Handle("/exec", s.Auth(s.execHandler)).Methods("POST")
}

func WriteJSON(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Println("Error marshaling json:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
