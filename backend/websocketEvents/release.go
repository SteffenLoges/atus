package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/bencode"
	"atus/backend/release"
	"atus/backend/sqlite"
	"atus/backend/websocket"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Release__GetDetails(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	releaseRow := sqlite.Conn.QueryRow(
		`SELECT 	
			uid, 
			hash,
			name, 
			name_raw,
			pre, 
			category,
			category_raw, 
			size, 
			added,
			source_uid,
			fileserver_uid,
			state,
			uploaded
		FROM releases 
		WHERE uid = ?
		LIMIT 1`,
		req.UID,
	)

	var uid, hash, name, nameRaw, pre, category, categoryRaw, addedRaw, sourceUID, fileserverUID, state string
	var uploaded sql.NullString
	var size int64

	err := releaseRow.Scan(&uid, &hash, &name, &nameRaw, &pre, &category, &categoryRaw, &size, &addedRaw, &sourceUID, &fileserverUID, &state, &uploaded)
	if err != nil {
		if err == sql.ErrNoRows {
			r.SetResponseCode(http.StatusNotFound)
			r.MarshalAndSendResponse("release not found")
			return
		}
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

	fileserverName := ""
	if fileserver := a.GetFileserverByUID(fileserverUID); fileserver != nil {
		fileserverName = fileserver.Fileserver.Name
	}

	downloadState, _ := a.GetDownloadState(fileserverUID, hash)
	metaFiles, _ := release.GetMetaFiles(uid, "")

	r.MarshalAndSendResponse(map[string]interface{}{
		"uid":            uid,
		"name":           name,
		"nameRaw":        nameRaw,
		"pre":            pre,
		"category":       category,
		"categoryRaw":    categoryRaw,
		"size":           size,
		"added":          added,
		"fileserverName": fileserverName,
		"sourceName":     sourceName,
		"downloadState":  downloadState,
		"metaFiles":      metaFiles,
		"state": map[string]interface{}{
			"state":      state,
			"uploadDate": uploaded.String,
		},
	})

}

func Release__GetFiles(r *websocket.Request) {

	var req struct {
		UID string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	releaseRow := sqlite.Conn.QueryRow(
		`SELECT 	
			uid, 
			source_uid
		FROM releases 
		WHERE uid = ?
		LIMIT 1`,
		req.UID,
	)

	var uid, sourceUID string

	err := releaseRow.Scan(&uid, &sourceUID)
	if err != nil {
		if err == sql.ErrNoRows {
			r.SetResponseCode(http.StatusNotFound)
			r.MarshalAndSendResponse("release not found")
			return
		}
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	meta, err := release.GetMetaFiles(req.UID, release.MetafileStateProcessed, release.MetafileTypeTorrent)
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	if len(meta) == 0 {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse("torrent file not found (01)")
		return
	}

	torrentFile, err := meta[0].GetFile()
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	dict, err := bencode.BDecode(torrentFile)
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(fmt.Sprintf("error decoding meta file: %s", err.Error()))
		return
	}

	fh := dict.GetFilesHierarchical()
	bencode.SortHierarchicalFiles(fh)

	r.MarshalAndSendResponse(fh)

}

func Release__Delete(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	row := sqlite.Conn.QueryRow("SELECT hash FROM releases	WHERE uid = ? LIMIT 1", req.UID)
	var hash string
	err := row.Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			r.SetResponseCode(http.StatusNotFound)
			r.MarshalAndSendResponse("release not found")
			return
		}
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	if err := a.DeleteRelease(req.UID, hash); err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(true)

}

func Release__Upload(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	rls := a.GetReleaseByUID(req.UID)
	if rls == nil {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse("release not found")
		return
	}

	if rls.State != release.StateDownloaded && rls.State != release.StateUploaded && rls.State != release.StateUploadError {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse("release is not in a state to be uploaded")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := a.UploadRelease(ctx, rls); err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(true)

}
