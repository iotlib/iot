package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/namsral/flag"

	"github.com/twinone/iot/backend/ws"
	"github.com/twinone/iot/backend/httpserver"
)

const (
	wsPath    = "/echo"
)

var addr = flag.String("addr", ":8080", "http service address and port")


func main() {
	flag.String(flag.DefaultConfigFlagname, "./config", "path to config file")
	flag.Parse()
	log.SetFlags(0)
	go ws.DefaultHub.Run()
	httpserver.Setup()

	http.HandleFunc(wsPath, ws.DefaultWSHandler)
	http.HandleFunc("/", httpserver.HomeHandler)
	http.HandleFunc("/auth/callback", httpserver.AuthHandler)

	fmt.Println("Listening at", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
