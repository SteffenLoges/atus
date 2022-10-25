package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/websocket"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

func Settings__SourcesManage_GetAll(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var ret []map[string]interface{}
	for _, source := range a.GetAllSources() {
		ret = append(ret, map[string]interface{}{
			"uid":                   source.Source.UID,
			"name":                  source.Source.Name,
			"favicon":               source.Source.Favicon,
			"enabled":               source.Source.Enabled,
			"timesChecked":          source.Source.TimesChecked,
			"lastChecked":           source.Source.LastCheck,
			"nextCheck":             source.Source.NextCheck,
			"sumTorrentsDownloaded": source.Source.SumTorrentsDownloaded,
			"sumImagesDownloaded":   source.Source.SumImagesDownloaded,
			"sumReleasesDownloaded": source.Source.SumReleasesDownloaded,
		})
	}

	// sort by name, uid ASC
	sort.Slice(ret, func(i, j int) bool {
		if ret[i]["name"].(string) == ret[j]["name"].(string) {
			return ret[i]["uid"].(string) < ret[j]["uid"].(string)
		}
		return ret[i]["name"].(string) < ret[j]["name"].(string)
	})

	r.MarshalAndSendResponse(ret)

}

func Settings__SourcesManage_Delete(r *websocket.Request) {

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

	if err := a.DeleteSource(s); err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(fmt.Sprintf("error deleting source: %s", err.Error()))
		return
	}

	r.MarshalAndSendResponse(true)
}

func Settings__SourcesManage_Toggle(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID   string
		Start bool
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

	if req.Start {
		go s.Enable(a, true)
	} else {
		go s.Disable()
	}

	r.MarshalAndSendResponse(true)

}
