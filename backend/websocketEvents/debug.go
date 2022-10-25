package websocketEvents

import (
	"atus/backend/logger"
	"atus/backend/websocket"
)

func Debug__GetCache(r *websocket.Request) {
	r.MarshalAndSendResponse(logger.GetCache())
}
