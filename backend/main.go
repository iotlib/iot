package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/namsral/flag"
	"github.com/twinone/iot/backend/db"
	"github.com/twinone/iot/backend/httpserver"
	"github.com/twinone/iot/backend/ws"
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

	sess := db.Init()
	defer sess.Close()

	hub := ws.DefaultHub
	go hub.Run()

	ss := httpserver.New(config, hub)

	r := mux.NewRouter()
	r.HandleFunc(wsPath, ws.GenWSHandler(hub))
	ss.RegisterHandlers(r)
	http.Handle("/", r)

	fmt.Println("Listening at", *config["addr"])
	log.Fatal(http.ListenAndServe(*config["addr"], nil))
}
