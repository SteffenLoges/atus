package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/websocket"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func Settings__SourcesEdit_Get(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	s := a.GetSourceByUID(req.UID)
	if s == nil {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse(fmt.Sprintf("source with UID %s not found", req.UID))
		return
	}

	var cookies []interface{}
	for _, c := range s.Source.Cookies {
		cookies = append(cookies, map[string]string{
			"name":  c.Name,
			"value": c.Value,
		})
	}

	r.MarshalAndSendResponse(map[string]interface{}{
		"name":              s.Source.Name,
		"favicon":           s.Source.Favicon,
		"rssURL":            s.Source.RSSURL.String(),
		"rssInterval":       s.Source.RSSInterval / time.Second,
		"requestWaitTime":   s.Source.RequestWaitTime / time.Millisecond,
		"cookies":           cookies,
		"metaPath":          s.Source.MetaPath,
		"metaPathUseAsKey":  s.Source.MetaPathUseAsKey,
		"imagePath":         s.Source.ImagePath,
		"imagePathUseAsKey": s.Source.ImagePathUseAsKey,
	})

}

func Settings__SourcesEdit_Save(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID               string
		Name              string
		RSSURL            string
		Cookies           []*http.Cookie
		RSSInterval       int64
		RequestWaitTime   int64
		MetaPath          string
		MetaPathUseAsKey  bool
		ImagePath         string
		ImagePathUseAsKey bool
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	s := a.GetSourceByUID(req.UID)
	if s == nil {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse(fmt.Sprintf("source with UID %s not found", req.UID))
		return
	}

	rssURL, err := url.Parse(req.RSSURL)
	if err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(fmt.Sprintf("rss url is not valid: %s", err))
		return
	}

	s.Source.Name = req.Name
	s.Source.RSSURL = rssURL
	s.Source.Cookies = req.Cookies
	s.Source.RSSInterval = time.Duration(req.RSSInterval) * time.Second
	s.Source.RequestWaitTime = time.Duration(req.RequestWaitTime) * time.Millisecond
	s.Source.MetaPath = req.MetaPath
	s.Source.MetaPathUseAsKey = req.MetaPathUseAsKey
	s.Source.ImagePath = req.ImagePath
	s.Source.ImagePathUseAsKey = req.ImagePathUseAsKey
	s.Source.Save()

	r.MarshalAndSendResponse(true)

}
