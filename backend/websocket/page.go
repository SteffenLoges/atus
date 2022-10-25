package websocket

import (
	"encoding/json"
	"net/http"
)

type Page struct {
	Name   string
	Params map[string]string
}

func currentPageEventHandler(r *Request) {
	if err := json.Unmarshal(r.Payload, &r.Client.CurrentPage); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}
}
