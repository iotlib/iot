package ws

import "log"

var DefaultHub = NewHub()

type Hub struct {
	register   chan *Conn
	unregister chan *Conn
	conns      map[*Conn]bool
	// Maps email to id
	OwnersToIds map[string][]string
	IdsToConns map[string]*Conn
}

func NewHub() *Hub {
	return &Hub{
		register: make(chan *Conn),
		unregister: make(chan *Conn),
		conns:   make(map[*Conn]bool),

		OwnersToIds: make(map[string][]string),
		IdsToConns: make(map[string]*Conn),
	}
}

func (h *Hub) Run() {
	defer func() {
		close(h.register)
		close(h.unregister)
		for c, _ := range h.conns {
			c.Close()
		}
	}()

	for {
		select {
		case conn := <-h.register:
			h.conns[conn] = true
			h.OwnersToIds[conn.owner] = append(h.OwnersToIds[conn.owner], conn.id)
			h.IdsToConns[conn.id] = conn

			log.Println("Registered conn")
		case conn := <-h.unregister:
			log.Println("Unregistered conn")
			delete(h.conns, conn)
		}
	}
}
