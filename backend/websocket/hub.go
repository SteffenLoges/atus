package websocket

import (
	"context"
	"encoding/json"
	"time"
)

type Hub struct {
	clients              map[*Client]bool
	broadcastChan        chan []byte
	register             chan *Client
	unregister           chan *Client
	eventHandlers        map[string]EventHandler
	Ctx                  context.Context
	OnClientConnected    func(*Client)
	OnClientDisconnected func(*Client)

	// Time allowed to write a message to the peer.
	WriteWait time.Duration

	// Time allowed to read the next pong message from the peer.
	PongWait time.Duration

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod time.Duration

	// Maximum message size in bytes allowed from peer.
	MaxMessageSize int64

	SendChanSize int64
}

func NewHub(ctx context.Context) *Hub {
	pongWait := 120 * time.Second

	b := &Hub{
		broadcastChan:        make(chan []byte),
		register:             make(chan *Client),
		unregister:           make(chan *Client),
		clients:              make(map[*Client]bool),
		eventHandlers:        make(map[string]EventHandler),
		Ctx:                  ctx,
		OnClientConnected:    func(c *Client) {},
		OnClientDisconnected: func(c *Client) {},

		WriteWait:      120 * time.Second,
		PongWait:       pongWait,
		PingPeriod:     (pongWait * 9) / 10,
		MaxMessageSize: 1024 * 5,
		SendChanSize:   1024 * 5,
	}

	b.SetEventHandler("CURRENT_PAGE", currentPageEventHandler)

	return b
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.OnClientConnected(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				client.onClientDisconnected()
			}
		case message := <-h.broadcastChan:
			for client := range h.clients {
				if client.Connected.IsNotSet() {
					continue
				}

				select {
				case client.sendChan <- message:
				default:
					client.onClientDisconnected()
				}
			}
		}
	}
}

func (h *Hub) GetAllClients() []*Client {
	var clients []*Client
	for c := range h.clients {
		if c.Connected.IsNotSet() {
			continue
		}

		clients = append(clients, c)
	}
	return clients
}

func (h *Hub) GetClientsByUID(uid string) []*Client {
	var clients []*Client
	for c := range h.clients {
		if c.Connected.IsSet() && c.UserUID == uid {
			clients = append(clients, c)
		}
	}
	return clients
}

// GetClientsByPage returns all clients that are currently on the given page with the given params.
func (h *Hub) GetClientsByPage(name string, query *map[string]string) []*Client {
	var clients []*Client
	for c := range h.clients {
		if c.Connected.IsNotSet() || c.CurrentPage == nil || c.CurrentPage.Name != name {
			continue
		}

		if query != nil {
			for k, v := range *query {
				if c.CurrentPage.Params[k] != v {
					continue
				}
			}
		}

		clients = append(clients, c)
	}
	return clients
}

func (h *Hub) Broadcast(action string, payload []byte) {
	msg := append([]byte(action), payloadSeperator...)
	msg = append(msg, payload...)
	go func(msg []byte) {
		h.broadcastChan <- msg
	}(msg)
}

func (h *Hub) MarshalAndBroadcast(action string, payload interface{}) {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	h.Broadcast(action, marshaled)
}
