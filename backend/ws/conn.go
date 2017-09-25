package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"log"
	"sync"
	"strings"
)

type State int

const (
	PendingHello State = iota
	PendingOwner
	Connected
)

type Command string

const (
	Hello = "HELLO"
	Owner = "OWNER"
	Bye   = "BYE"
)

const (
	// Number of messages in receiving and sending queues before blocking
	queueSize = 16

	// Time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 7) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// WSHandler using the DefaultHub
var DefaultWSHandler = GenWSHandler(DefaultHub)


var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Conn struct {
	hub  *Hub
	ws   *websocket.Conn
	Send chan []byte
	Recv chan []byte

	State State
	owner string
	id    string

	mx     sync.Mutex
	closed bool
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if (!ok) {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			c.ws.WriteMessage(websocket.TextMessage, msg)
		case <-ticker.C:
			log.Print("ping... ")
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Conn) readPump() {
	defer c.Close()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			//if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			//	log.Printf("error: %v", err)
			//}
			break
		}
		log.Println("RECV:", string(message))
		c.processMessage(message)

		c.mx.Lock()
		if !c.closed {
			c.Recv <- message
		}
		c.mx.Unlock()
	}
}

// Returns true if message is ok (expected message)
func (c *Conn) processMessage(message []byte) bool {
	msg := string(message)
	ss := strings.SplitN(msg, " ", 2)
	cmd := strings.Trim(ss[0], " \t\r\n")

	switch cmd {
	case Hello:
		if len(ss) < 2 || c.State != PendingHello {
			c.Close();
			return false
		}
		c.id = ss[1]
		// TODO check if id is ok
		c.State = PendingOwner
	case Owner:
		if len(ss) < 2 || c.State != PendingOwner {
			c.Close()
			return false
		}
		c.owner = ss[1]
		// TODO check if owner is registered etc
		c.State = Connected
		DefaultHub.register <- c
	case Bye:
		c.Close()
	default:
		if c.State != Connected {
			log.Println("Unexpected msg:", msg)
			c.Close()
			return false
		}
	}
	return true
}

func (c *Conn) Close() {
	c.mx.Lock()
	defer c.mx.Unlock()

	if !c.closed {
		c.closed = true
		close(c.Send)
		close(c.Recv)
		c.ws.Close()
		if c.State == Connected {
			c.hub.unregister <- c
		}
	}
}


// Generate a new WS Handler associated to a Hub
func GenWSHandler(hub *Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		conn := &Conn{
			Send:  make(chan []byte, queueSize),
			Recv:  make(chan []byte, queueSize),
			State: PendingHello,
			ws:    ws,
			hub:   hub,
		}

		go conn.writePump()
		go conn.readPump()

	}
}
