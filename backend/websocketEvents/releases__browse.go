package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/release"
	"atus/backend/sqlite"
	"atus/backend/websocket"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func Releases__Browse_Get(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		Name     string
		Category string
		State    string
		Offset   int
		Limit    int
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	// --------------------------------------------

	var clauseArr []string
	var bindings []interface{}

	if req.Name != "" {
		clauseArr = append(clauseArr, "name LIKE ?")
		bindings = append(bindings, "%"+req.Name+"%")
	}

	if req.Category != "" && req.Category != "all" {
		clauseArr = append(clauseArr, "category = ?")
		bindings = append(bindings, strings.ToUpper(req.Category))
	}

	if req.State != "" && req.State != "all" {
		if req.State == "all_but_uploaded" {
			clauseArr = append(clauseArr, "state != ?")
			bindings = append(bindings, release.StateUploaded)
		} else {
			clauseArr = append(clauseArr, "state = ?")
			bindings = append(bindings, strings.ToUpper(req.State))
		}
	}

	var clause string
	if len(clauseArr) > 0 {
		clause += " WHERE " + strings.Join(clauseArr, " AND ")
	}

	// --------------------------------------------

	var count int
	err := sqlite.Conn.QueryRow("SELECT COUNT(1) FROM releases"+clause, bindings...).Scan(&count)
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	// --------------------------------------------

	mainQuery :=
		`SELECT 
			uid, 
			hash,
			name, 
			pre, 
			state,
			category, 
			category_raw,
			size,
			added,
			source_uid,
			fileserver_uid,
			uploaded
		FROM releases`

	mainQuery += clause + " ORDER BY added DESC"

	if req.Limit > 0 {
		mainQuery += " LIMIT " + fmt.Sprintf("%d", req.Limit) + " OFFSET " + fmt.Sprintf("%d", req.Offset)
	}

	releaseRows, err := sqlite.Conn.Query(mainQuery, bindings...)
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	defer releaseRows.Close()

	var releases []map[string]interface{}
	for releaseRows.Next() {
		var uid, name, pre, category, categoryRaw, addedRaw, hash, sourceUID, fileserverUID string
		var uploaded sql.NullString
		var state release.ReleaseState
		var size int64

		if err := releaseRows.Scan(
			&uid,
			&hash,
			&name,
			&pre,
			&state,
			&category,
			&categoryRaw,
			&size,
			&addedRaw,
			&sourceUID,
			&fileserverUID,
			&uploaded,
		); err != nil {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse(err.Error())
			return
		}

		added, err := time.Parse(time.RFC3339, addedRaw)
		if err != nil {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse(err.Error())
			return
		}

		sourceName := ""
		if source := a.GetSourceByUID(sourceUID); source != nil {
			sourceName = source.Source.Name
		}

		downloadState, _ := a.GetDownloadState(fileserverUID, hash)

		release := map[string]interface{}{
			"uid":           uid,
			"hash":          hash,
			"name":          name,
			"pre":           pre,
			"category":      category,
			"categoryRaw":   categoryRaw,
			"size":          size,
			"added":         added,
			"fileserverUID": fileserverUID,
			"sourceName":    sourceName,
			"downloadState": downloadState,
			"state": map[string]interface{}{
				"state":      state,
				"uploadDate": uploaded.String,
			},
		}

		releases = append(releases, release)
	}

	for _, r := range releases {
		metaFiles, _ := release.GetMetaFiles(r["uid"].(string), "")
		r["metaFiles"] = metaFiles
	}

	r.MarshalAndSendResponse(map[string]interface{}{
		"count":    count,
		"releases": releases,
	})

}
