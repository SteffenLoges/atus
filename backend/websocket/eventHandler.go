package websocket

import (
	"net/http"
)

type EventHandler func(r *Request)

func (h *Hub) SetEventHandler(path string, handler EventHandler) {
	h.eventHandlers[path] = handler
}

func (h *Hub) getEventHandler(path string) EventHandler {

	if handler, ok := h.eventHandlers[path]; ok {
		return handler
	}

	// Send 404 if action handler was not found
	return func(r *Request) {
		r.SetResponseCode(http.StatusNotFound)
		r.SendResponse([]byte(http.StatusText(http.StatusNotFound)))
	}

}
