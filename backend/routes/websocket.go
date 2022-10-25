package routes

import (
	"atus/backend/atus"
	"atus/backend/user"
	"atus/backend/websocket"
	"net/http"
)

func SocketUserHandler(h *websocket.Hub, a *atus.ATUS) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.ContextKey).(*user.User)
		h.Serve(GetIP(r), u.UID, w, r)
	}
}
