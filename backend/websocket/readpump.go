package websocket

import (
	"atus/backend/logger"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.OnClientDisconnected(c)
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(h.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(h.PongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(h.PongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handleMessage(h, string(message))
	}
}

var (
	messageRegExp, _            = regexp.Compile(`^(?P<action>[A-Za-z0-9_-]+)(:(?P<requestID>[A-Za-z0-9]{15}))?(:(?P<payload>.+))?$`)
	messageRegExpActionIndex    = messageRegExp.SubexpIndex("action")
	messageRegExpRequestIDIndex = messageRegExp.SubexpIndex("requestID")
	messageRegExpPayloadIndex   = messageRegExp.SubexpIndex("payload")
)

func (c *Client) handleMessage(h *Hub, message string) {

	matches := messageRegExp.FindStringSubmatch(message)

	if matches == nil {
		logger.ForceConsole().Debugf("Error while reading websocket message. UID: %s Message: %s", c.UserUID, message)
		return
	}

	action := matches[messageRegExpActionIndex]
	requestID := matches[messageRegExpRequestIDIndex]
	payload := matches[messageRegExpPayloadIndex]

	handler := h.getEventHandler(action)

	handler(&Request{
		Hub:          h,
		Client:       c,
		action:       action,
		Payload:      []byte(payload),
		requestID:    requestID,
		responseCode: http.StatusOK,
	})

}

type Request struct {
	Hub          *Hub
	Client       *Client
	Payload      []byte
	requestID    string
	action       string
	responseCode int
}

func (r *Request) SetResponseCode(code int) {
	r.responseCode = code
}

func (r *Request) SendResponse(payload []byte) {
	r.Client.SendResponse(r.responseCode, r.requestID, r.action, payload)
}

func (r *Request) MarshalAndSendResponse(payload interface{}) {
	r.Client.MarshalAndSendResponse(r.responseCode, r.requestID, r.action, payload)
}
