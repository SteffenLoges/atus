package websocketEvents

import (
	"atus/backend/config"
	"atus/backend/helpers"
	"atus/backend/logger"
	"atus/backend/websocket"
)

func resetAPIToken() string {
	t := helpers.GetUUID()
	config.Set("API__AUTH_TOKEN", t)
	logger.Info("new API token generated")
	return t
}

func Settings__APIManage_ResetAPIToken(r *websocket.Request) {
	r.MarshalAndSendResponse(map[string]interface{}{
		"authToken": resetAPIToken(),
	})
}

func Settings__APIManage_GetToken(r *websocket.Request) {
	t := config.GetString("API__AUTH_TOKEN")
	if t == "" {
		t = resetAPIToken()
	}

	r.MarshalAndSendResponse(map[string]interface{}{
		"authToken": t,
	})
}
