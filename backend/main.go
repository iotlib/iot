package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"fmt"
)

const Path = "/echo"

var addr = flag.String("addr", ":8080", "http service address")

var conns map[*websocket.Conn]bool

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func broadcast(msg string) {
	for conn := range (conns) {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	log.Println("Client connected")

	conns[conn] = true

	// reader
	go func() {
		defer func() {
			conn.Close()
			delete(conns, conn)
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			log.Printf("recv: %s", message)
			broadcast(string(message))
		}
	}()
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, Path)
}

func main() {
	conns = make(map[*websocket.Conn]bool)

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc(Path, echo)
	http.HandleFunc("/", home)
	fmt.Println("Listening at", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
		var loc = window.location, new_uri;
		if (loc.protocol === "https:")
			new_uri = "wss:";
		else new_uri = "ws:";

		new_uri += "//" + loc.host;
		new_uri += {{.}};
        ws = new WebSocket(new_uri);
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: CMD " + input.value);
        ws.send("CMD " + input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
