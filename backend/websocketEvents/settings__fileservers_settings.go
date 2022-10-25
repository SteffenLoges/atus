package websocketEvents

import (
	"atus/backend/config"
	"atus/backend/websocket"
	"encoding/json"
	"net/http"
)

func Settings__FileserversSettings_Get(r *websocket.Request) {
	r.MarshalAndSendResponse(map[string]interface{}{
		"allocationMethod": config.GetString("FILESERVER__ALLOCATION_METHOD"),
		"downloadLabel":    config.GetString("FILESERVER__DOWNLOAD_LABEL"),
		"uploadLabel":      config.GetString("FILESERVER__UPLOAD_LABEL"),
	})
}

func Settings__FileserversSettings_Save(r *websocket.Request) {

	var req struct {
		AllocationMethod string
		DownloadLabel    string
		UploadLabel      string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	config.Set("FILESERVER__ALLOCATION_METHOD", req.AllocationMethod)
	config.Set("FILESERVER__DOWNLOAD_LABEL", req.DownloadLabel)
	config.Set("FILESERVER__UPLOAD_LABEL", req.UploadLabel)

	r.MarshalAndSendResponse(true)

}
