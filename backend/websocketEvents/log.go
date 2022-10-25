package websocketEvents

import (
	"atus/backend/sqlite"
	"atus/backend/websocket"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func Log__Get(r *websocket.Request) {

	var req struct {
		ReleaseUID string
		Type       string
		Severity   int
		Offset     int
		Limit      int
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	// --------------------------------------------

	var clauseArr []string
	var bindings []interface{}

	if req.ReleaseUID != "" {
		clauseArr = append(clauseArr, "release_uid = ?")
		bindings = append(bindings, req.ReleaseUID)
	}

	if req.Type != "" {
		clauseArr = append(clauseArr, "type = ?")
		bindings = append(bindings, req.Type)
	}

	if req.Severity != -1 {
		clauseArr = append(clauseArr, "severity = ?")
		bindings = append(bindings, req.Severity)
	}

	var clause string
	if len(clauseArr) > 0 {
		clause += " WHERE " + strings.Join(clauseArr, " AND ")
	}

	// --------------------------------------------

	var count int
	err := sqlite.Conn.QueryRow("SELECT COUNT(1) FROM log"+clause, bindings...).Scan(&count)
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	// --------------------------------------------

	var types []string
	typesRow, err := sqlite.Conn.Query("SELECT DISTINCT type FROM log")
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	defer typesRow.Close()

	for typesRow.Next() {
		var t string
		err = typesRow.Scan(&t)
		if err != nil {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse(err.Error())
			return
		}

		types = append(types, t)
	}

	// --------------------------------------------

	mainQuery :=
		`SELECT 
			severity,
			type,
			added,
			message,
			release_uid,
			(SELECT name FROM releases WHERE uid = release_uid)
		FROM log`

	mainQuery += clause + " ORDER BY added DESC"

	if req.Limit > 0 {
		mainQuery += " LIMIT " + fmt.Sprintf("%d", req.Limit) + " OFFSET " + fmt.Sprintf("%d", req.Offset)
	}

	logRows, err := sqlite.Conn.Query(mainQuery, bindings...)
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	defer logRows.Close()

	var entries []map[string]interface{}
	for logRows.Next() {
		var severity int
		var typeStr, added, message, releaseUID string
		var releaseName sql.NullString

		if err := logRows.Scan(&severity, &typeStr, &added, &message, &releaseUID, &releaseName); err != nil {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse(err.Error())
			return
		}

		entries = append(entries, map[string]interface{}{
			"severity":    severity,
			"type":        typeStr,
			"added":       added,
			"message":     message,
			"releaseUID":  releaseUID,
			"releaseName": releaseName.String,
		})

	}

	r.MarshalAndSendResponse(map[string]interface{}{
		"count":   count,
		"entries": entries,
		"types":   types,
	})

}

func Log__Clear(r *websocket.Request) {

	_, err := sqlite.Conn.Exec("DELETE FROM log")

	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(true)

}
