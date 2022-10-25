package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/category"
	"atus/backend/config"
	"atus/backend/helpers"
	"atus/backend/websocket"
	"encoding/json"
	"net/http"
	"sort"
)

func Settings__FiltersMisc_GetAll(r *websocket.Request) {
	r.MarshalAndSendResponse(map[string]interface{}{
		"maxAge": config.GetInt64("FILTERS__MAX_AGE"),
	})
}

func Settings__FiltersMisc_Save(r *websocket.Request) {

	var req struct {
		MaxAge int64 `json:"maxAge"`
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	config.Set("FILTERS__MAX_AGE", req.MaxAge)

	r.MarshalAndSendResponse(true)

}

func Settings__FiltersCategories_GetAll(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var ret []map[string]interface{}
	for _, c := range a.GetAllCategories() {
		ret = append(ret, map[string]interface{}{
			"enabled":  c.Enabled,
			"name":     c.Name,
			"includes": c.Includes,
			"excludes": c.Excludes,
			"maxSize":  c.MaxSize / helpers.GiB,
		})
	}

	// sort by name
	sort.Slice(ret, func(i, j int) bool {
		return ret[i]["name"].(category.Name) < ret[j]["name"].(category.Name)
	})

	r.MarshalAndSendResponse(ret)

}

func Settings__FiltersCategories_Save(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		Categories []*category.Category
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	for _, c := range req.Categories {
		c.MaxSize = c.MaxSize * helpers.GiB
		if err := a.UpdateCategory(c); err != nil {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse(err.Error())
			return
		}
	}

	r.MarshalAndSendResponse(true)

}
