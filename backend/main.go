package main

import (
	"flag"
	"log"
	"net/http"
	"fmt"

	"github.com/twinone/iot/backend/ws"
	"github.com/twinone/iot/backend/httpserver"
)


const (
	wsPath = "/echo"
	indexPath = "/"
)

var addr = flag.String("addr", ":8080", "http service address and port")


func main() {
	flag.Parse()
	log.SetFlags(0)
	go ws.DefaultHub.Run()

	http.HandleFunc(wsPath, ws.DefaultWSHandler)
	http.HandleFunc(indexPath, httpserver.GenHomeHandler(wsPath))

	fmt.Println("Listening at", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

