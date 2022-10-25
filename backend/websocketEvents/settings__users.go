package websocketEvents

import (
	"atus/backend/user"
	"atus/backend/websocket"
	"encoding/json"
	"net/http"
)

func Settings__Users_Get(r *websocket.Request) {

	var req struct {
		UID string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	u, err := user.GetByUID(req.UID)
	if err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(u)
}

func Settings__Users_GetAll(r *websocket.Request) {

	users, err := user.GetAll()
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(users)
}

func Settings__Users_Add(r *websocket.Request) {

	var req struct {
		Username string
		Password string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	if _, err := user.Add(req.Username, req.Password, false); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(true)

}

func Settings__Users_Update(r *websocket.Request) {

	var req struct {
		UID      string
		Username string
		Password string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	u, err := user.GetByUID(req.UID)
	if err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	if req.Username != "" && u.Name != req.Username {
		u.Name = req.Username
		if err := u.Save(); err != nil {
			r.SetResponseCode(http.StatusBadRequest)
			r.MarshalAndSendResponse(err.Error())
			return
		}
	}

	if req.Password != "" {
		if err := u.SetPassword(req.Password); err != nil {
			r.SetResponseCode(http.StatusBadRequest)
			r.MarshalAndSendResponse(err.Error())
			return
		}
	}

	r.MarshalAndSendResponse(true)

}

func Settings__Users_Delete(r *websocket.Request) {

	var req struct {
		UID string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	if err := user.Delete(req.UID); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(true)

}
