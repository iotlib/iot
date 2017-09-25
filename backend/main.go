package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/namsral/flag"
	"github.com/twinone/iot/backend/ws"
	"github.com/twinone/iot/backend/httpserver"
	"github.com/gorilla/mux"
)

const (
	wsPath = "/echo"
)

var config map[string]*string

func main() {

	config = map[string]*string{
		"addr":                flag.String("addr", ":8080", "http service address and port"),
		"callback_url":        flag.String("callback_url", "", "OAuth Callback URL"),
		"client_id":           flag.String("client_id", "", "OAuth Client ID"),
		"client_secret":       flag.String("client_secret", "", "OAuth Client Secret"),
		"cookie_store_secret": flag.String("cookie_store_secret", "", "Cookie-store secret"),
	}
	flag.String(flag.DefaultConfigFlagname, "./config", "path to config file")

	flag.Parse()
	log.SetFlags(0)

	hub := ws.DefaultHub
	go hub.Run()

	ss := httpserver.New(config, hub)

	r := mux.NewRouter()
	ss.RegisterHandlers(r)
	r.HandleFunc(wsPath, ws.GenWSHandler(hub))
	http.Handle("/", r)

	fmt.Println("Listening at", *config["addr"])
	log.Fatal(http.ListenAndServe(*config["addr"], nil))
}
