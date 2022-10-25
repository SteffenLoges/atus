package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/websocket"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

func Settings__FileserversManage_GetAll(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var ret []map[string]interface{}
	for _, fs := range a.GetAllFileservers() {

		r := map[string]interface{}{
			"uid":             fs.Fileserver.UID,
			"name":            fs.Fileserver.Name,
			"enabled":         fs.Fileserver.Enabled,
			"filesDownloaded": fs.Fileserver.SumFilesDownloaded,
		}

		if fs.Fileserver.Statistics != nil {
			r["serverLoad"] = fs.Fileserver.Statistics.ServerLoad
			r["diskFreeSpace"] = fs.Fileserver.Statistics.DiskFreeSpace
			r["diskTotalSpace"] = fs.Fileserver.Statistics.DiskTotalSpace
		}

		ret = append(ret, r)
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

func Settings__FileserversManage_Delete(r *websocket.Request) {

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

	if err := a.DeleteFileserver(s); err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(fmt.Sprintf("error deleting fileserver: %s", err.Error()))
		return
	}

	r.MarshalAndSendResponse(true)

}

func Settings__FileserversManage_Toggle(r *websocket.Request) {

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

	s := a.GetFileserverByUID(req.UID)
	if s == nil {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse(fmt.Sprintf("fileserver with UID %s not found", req.UID))
		return
	}

	if req.Start {
		go s.Enable(a)
	} else {
		go s.Disable(a)
	}

	r.MarshalAndSendResponse(true)

}
