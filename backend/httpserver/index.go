package httpserver

import "net/http"

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	c := s.GetCookie(r)
	u := s.GetUser(c)

	if u == nil {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)

	}
}
