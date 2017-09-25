package ws

import (
	"log"
	"github.com/twinone/iot/backend/model"
)

var DefaultHub = NewHub()

type Hub struct {
	register   chan *Conn
	unregister chan *Conn
	conns      map[*Conn]bool
	// Maps email to id
	OwnersToIds map[string]map[string]bool
	IdsToConns  map[string]*Conn
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Conn),
		unregister: make(chan *Conn),
		conns:      make(map[*Conn]bool),

		OwnersToIds: make(map[string]map[string]bool),
		IdsToConns:  make(map[string]*Conn),
	}
}

func (h *Hub) GetDevices(owner string) []*model.Device {
	idmap := h.OwnersToIds[owner]
	res := make([]*model.Device, 0, len(idmap))

	for k, _ := range idmap {

		log.Println(k)
		if conn, ok := h.IdsToConns[k]; ok {
			res = append(res, conn.device)
		} else {
			// device is offline
			// TODO
		}
	}
	return res
}

func (h *Hub) Run() {
	defer func() {

		for c, _ := range h.conns {
			c.Close()
		}
		close(h.register)
		close(h.unregister)
	}()

	for {
		select {
		case conn := <-h.register:
			h.conns[conn] = true
			if _, ok := h.OwnersToIds[conn.device.Owner]; !ok {
				h.OwnersToIds[conn.device.Owner] = make(map[string]bool)
			}
			h.OwnersToIds[conn.device.Owner][conn.device.Id] = true
			h.IdsToConns[conn.device.Id] = conn

			log.Println("Registered conn")
		case conn := <-h.unregister:
			log.Println("Unregistered conn")
			delete(h.conns, conn)
		}
	}
}
