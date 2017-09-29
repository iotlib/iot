package ws

import (
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

func (h *Hub) GetConns(owner string) []*Conn {
	idmap := h.OwnersToIds[owner]
	res := make([]*Conn, 0, len(idmap))

	for k, _ := range idmap {
		if conn, ok := h.IdsToConns[k]; ok {
			res = append(res, conn)
		}
	}
	return res
}

func (h *Hub) GetDevices(owner string) []*model.Device {
	idmap := h.OwnersToIds[owner]
	res := make([]*model.Device, 0, len(idmap))

	for k, _ := range idmap {
		if conn, ok := h.IdsToConns[k]; ok {
			res = append(res, conn.Device)
		}
	}
	return res
}

func (h *Hub) Run() {
	cleanup := func(conn *Conn) {
		delete(h.IdsToConns, conn.Device.Id)
		delete(h.conns, conn)
		conn.Close()
	}

	defer func() {

		for c, _ := range h.conns {
			cleanup(c)
		}
		close(h.register)
		close(h.unregister)
	}()

	for {
		select {
		case conn := <-h.register:
			h.conns[conn] = true
			if _, ok := h.OwnersToIds[conn.Device.Owner]; !ok {
				h.OwnersToIds[conn.Device.Owner] = make(map[string]bool)
			}
			h.OwnersToIds[conn.Device.Owner][conn.Device.Id] = true
			h.IdsToConns[conn.Device.Id] = conn

			//log.Println("Registered conn")
		case conn := <-h.unregister:
			//log.Println("Unregistered conn")
			cleanup(conn)
		}
	}
}
