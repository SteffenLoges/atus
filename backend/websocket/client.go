package websocket

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"atus/backend/helpers"

	"github.com/gorilla/websocket"
	"github.com/tevino/abool/v2"
)

var (
	payloadSeperator = []byte{':'}
	messageSeperator = []byte{'\n'}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	sendChan    chan []byte
	Connected   *abool.AtomicBool
	CurrentPage *Page
	UserUID     string
	IP          string
}

var ClientContextKey helpers.ContextKey = "client"

func (h *Hub) Serve(ip, userUID string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	client := &Client{
		hub:       h,
		conn:      conn,
		sendChan:  make(chan []byte, h.SendChanSize),
		Connected: abool.NewBool(true),
		UserUID:   userUID,
		IP:        ip,
	}

	h.register <- client

	go client.writePump(h)
	go client.readPump(h)
}

func (c *Client) onClientDisconnected() {
	c.Connected.UnSet()
	close(c.sendChan)
	delete(c.hub.clients, c)
}

func (c *Client) Send(action string, payload []byte) {
	if c.Connected.IsNotSet() {
		return
	}

	msg := [][]byte{
		[]byte(action),
		payload,
	}
	c.sendChan <- bytes.Join(msg, payloadSeperator)
}

func (c *Client) SendResponse(statusCode int, requestID, action string, payload []byte) {
	if c.Connected.IsNotSet() {
		return
	}

	msg := [][]byte{
		[]byte(action),
		[]byte(requestID),
		[]byte(strconv.Itoa(statusCode)),
	}

	if len(payload) != 0 {
		msg = append(msg, payload)
	}

	c.sendChan <- bytes.Join(msg, payloadSeperator)
}

func (c *Client) MarshalAndSend(action string, payload interface{}) {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	c.Send(action, marshaled)
}

func (c *Client) MarshalAndSendResponse(statusCode int, requestID, action string, payload interface{}) {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	c.SendResponse(statusCode, requestID, action, marshaled)
}
