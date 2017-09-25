package httpserver

import (
	"net/http"
	"html/template"
)

var homeTemplate = template.Must(template.ParseFiles("www/index.tmpl"))

func GenHomeHandler(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		homeTemplate.Execute(w, path)
	}
}
