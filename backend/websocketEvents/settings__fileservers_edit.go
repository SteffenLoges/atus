package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/helpers"
	"atus/backend/websocket"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func Settings__FileserversEdit_Get(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	s := a.GetFileserverByUID(req.UID)
	if s == nil {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse(fmt.Sprintf("fileserver with UID %s not found", req.UID))
		return
	}

	ret := map[string]interface{}{
		"uid":                s.Fileserver.UID,
		"name":               s.Fileserver.Name,
		"listInterval":       s.Fileserver.ListInterval / time.Second,
		"minFreeDiskSpace":   s.Fileserver.MinFreeDiskSpace / helpers.GiB,
		"statisticsInterval": s.Fileserver.StatisticsInterval / time.Second,
		"url":                s.Fileserver.URL.String(),
	}

	if s.Fileserver.Statistics != nil {
		ret["diskFreeSpace"] = s.Fileserver.Statistics.DiskFreeSpace
		ret["diskTotalSpace"] = s.Fileserver.Statistics.DiskTotalSpace
	}

	r.MarshalAndSendResponse(ret)

}

func Settings__FileserversEdit_Save(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID                string
		Name               string
		URL                string
		ListInterval       int64
		StatisticsInterval int64
		MinFreeDiskSpace   int64
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	s := a.GetFileserverByUID(req.UID)
	if s == nil {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse(fmt.Sprintf("fileserver with UID %s not found", req.UID))
		return
	}

	urlParsed, err := url.Parse(req.URL)
	if err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	listInterval := time.Duration(req.ListInterval) * time.Second
	if listInterval > time.Second {
		s.Fileserver.ListInterval = listInterval
	}

	statisticsInterval := time.Duration(req.StatisticsInterval) * time.Second
	if statisticsInterval > time.Second {
		s.Fileserver.StatisticsInterval = statisticsInterval
	}

	if req.MinFreeDiskSpace > 0 {
		s.Fileserver.MinFreeDiskSpace = helpers.GiB * req.MinFreeDiskSpace
	}

	s.Fileserver.Name = req.Name
	if s.Fileserver.Name == "" {
		s.Fileserver.Name = urlParsed.Host
	}

	s.Fileserver.URL = urlParsed

	if err := s.Save(); err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(true)

}
