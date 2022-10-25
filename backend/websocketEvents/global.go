package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/websocket"
)

func Global_ForceUpdateFileserverStatistics(r *websocket.Request) {
	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)
	r.Client.MarshalAndSend("FILESERVER_STATISTICS", a.GetFileserverStatistics())
}
