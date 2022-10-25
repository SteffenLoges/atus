package websocketEvents

import (
	"atus/backend/config"
	"atus/backend/helpers"
	"atus/backend/websocket"
	"encoding/json"
	"net/http"
)

func Settings__SamplesManage_GetAll(r *websocket.Request) {
	r.MarshalAndSendResponse(map[string]interface{}{
		"enabled":        config.GetBool("SAMPLES__ENABLED"),
		"sumScreenshots": config.GetInt64("SAMPLES__SUM_SCREENSHOTS"),
		"minSize":        config.GetInt64("SAMPLES__MIN_SIZE") / helpers.MiB,
		"maxSize":        config.GetInt64("SAMPLES__MAX_SIZE") / helpers.MiB,
	})
}

func Settings__SamplesManage_Save(r *websocket.Request) {

	var req struct {
		Enabled        bool
		SumScreenshots int64
		MinSize        int64
		MaxSize        int64
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	config.Set("SAMPLES__ENABLED", req.Enabled)
	config.Set("SAMPLES__SUM_SCREENSHOTS", req.SumScreenshots)
	config.Set("SAMPLES__MIN_SIZE", req.MinSize*helpers.MiB)
	config.Set("SAMPLES__MAX_SIZE", req.MaxSize*helpers.MiB)

	r.MarshalAndSendResponse(true)

}
