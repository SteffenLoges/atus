package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/config"
	"atus/backend/helpers"
	"atus/backend/release"
	"atus/backend/websocket"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func Settings__Upload_GetAll(r *websocket.Request) {
	r.MarshalAndSendResponse(map[string]interface{}{
		"uploadUserID":             config.GetString("UPLOAD__USER_ID"),
		"uploadUserAnnounceURL":    config.GetString("UPLOAD__USER_ANNOUNCE_URL"),
		"uploadTrackerAnnounceURL": config.GetString("UPLOAD__TRACKER_ANNOUNCE_URL"),
		"uploadAPIURL":             config.GetString("UPLOAD__API_URL"),
		"uploadCreatedBy":          config.GetString("UPLOAD__CREATED_BY"),
		"uploadComment":            config.GetString("UPLOAD__COMMENT"),
	})
}

func Settings__Upload_Save(r *websocket.Request) {

	var req struct {
		UploadUserID             string
		UploadUserAnnounceURL    string
		UploadTrackerAnnounceURL string
		UploadAPIURL             string
		UploadCreatedBy          string
		UploadComment            string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	config.Set("UPLOAD__USER_ID", req.UploadUserID)
	config.Set("UPLOAD__USER_ANNOUNCE_URL", req.UploadUserAnnounceURL)
	config.Set("UPLOAD__TRACKER_ANNOUNCE_URL", req.UploadTrackerAnnounceURL)
	config.Set("UPLOAD__API_URL", req.UploadAPIURL)
	config.Set("UPLOAD__CREATED_BY", req.UploadCreatedBy)
	config.Set("UPLOAD__COMMENT", req.UploadComment)

	atus.SetSetupStepDone(r.Hub, atus.SetupStepUploadConfigured)

	r.MarshalAndSendResponse(true)

}

func Settings__Upload_UploadTestTorrent(r *websocket.Request) {

	var req struct {
		UploadUserID             string
		UploadUserAnnounceURL    string
		UploadTrackerAnnounceURL string
		UploadAPIURL             string
		UploadCreatedBy          string
		UploadComment            string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	uid := "TESTFILE"

	// check if folder exists
	if _, err := os.Stat(path.Join(config.Base.Folders.Data, uid)); os.IsNotExist(err) {
		err := os.MkdirAll(path.Join(config.Base.Folders.Data, uid), 0755)
		if err != nil {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse(err.Error())
			return
		}
	}

	// create test images
	images := []release.MetaFileType{
		release.MetafileTypeImage,
		release.MetafileTypeProofImage,
		release.MetafileTypeScreenImage,
		release.MetafileTypeSourceImage,
	}

	sumSampleImages := config.GetInt64("SAMPLES__SUM_SCREENSHOTS")
	if sumSampleImages == 0 {
		sumSampleImages = 3
	}
	for i := int64(1); i <= sumSampleImages; i++ {
		images = append(images, release.MetafileTypeScreenImageFromSample)
	}

	var metaFiles []*release.MetaFile

	for i, imageType := range images {
		filename := fmt.Sprintf("%s_%d.jpg", imageType, i)
		savePath := path.Join(config.Base.Folders.Data, uid, filename)
		if _, err := os.Stat(savePath); errors.Is(err, os.ErrNotExist) {

			img := helpers.CreateSquareImage(400)
			img.AddLabel("ATUS", color.RGBA{25, 25, 25, 255}, 0)
			img.AddLabel("TEST "+string(imageType), color.RGBA{100, 25, 25, 180}, 50)

			if err := img.Save(savePath); err != nil {
				r.SetResponseCode(http.StatusInternalServerError)
				r.MarshalAndSendResponse(err.Error())
				return
			}
		}

		metaFiles = append(metaFiles, &release.MetaFile{
			ReleaseUID: uid,
			FileName:   filename,
			Type:       imageType,
			State:      release.MetafileStateProcessed,
			Info: release.MetaInfo{
				"type": string(imageType),
			},
		})
	}

	// create info file
	// this file serves no purpose but to inform the user that this is a test release
	if _, err := os.Stat(path.Join(config.Base.Folders.Data, uid, "info.txt")); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path.Join(config.Base.Folders.Data, uid, "info.txt"))
		if err != nil {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse(err.Error())
			return
		}
		f.WriteString("Files in this folder are used to test upload functionality.")
		f.WriteString("\n")
		f.WriteString("You can safely delete this folder.\n")
		f.Close()
	}

	torrent := []byte(
		"d10:created by13:uTorrent/160013:creation datei1662738966e8:encoding5:UTF-84:" +
			"infod5:filesld6:lengthi1e4:pathl14:testfile-1.txteed6:lengthi1e4:pathl14:testfile-2.txteee4:name4" +
			":test12:piece lengthi65536e6:pieces20:\xdd\xfe\x163E\xd38\x19:½\xc1\x83\xf8\xe9\xdc\xff\x90KCee",
	)

	nfo := []string{
		"             _______ _________          _______ ",
		"            (  ___  )\\__   __/|\\     /|(  ____ \\",
		"            | (   ) |   ) (   | )   ( || (    \\/",
		"            | (___) |   | |   | |   | || (_____ ",
		"            |  ___  |   | |   | |   | |(_____  )",
		"            | (   ) |   | |   | |   | |      ) |",
		"            | )   ( |   | |   | (___) |/\\____) |",
		"            |/     \\|   )_(   (_______)\\_______)",
		"              Automatic Torrent Upload Script",
		"",
		"            https://github.com/SteffenLoges/atus",
		"",
		"",
		"",
		"Check if the torrent is valid by downloading it from the tracker",
		"and adding it to a torrent client.",
		"",
		"If something is wrong, delete the torrent and try again.",
		"",
		"Once everything is setup to your liking, you can proceed to the",
		"next step.",
	}

	rls := &atus.Release{
		UID:       uid,
		Name:      "ATUS Test Release",
		Pre:       time.Now(),
		MetaFiles: metaFiles,
	}

	d := &atus.Destination{
		TrackerAnnounceURL: req.UploadTrackerAnnounceURL,
		Comment:            req.UploadComment,
		CreatedBy:          req.UploadCreatedBy,
		UserID:             req.UploadUserID,
		APIURL:             req.UploadAPIURL,
		APIAuthToken:       config.GetString("API__AUTH_TOKEN"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := d.UploadRelease(ctx, rls, torrent, []byte(strings.Join(nfo, "\n")))
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	r.MarshalAndSendResponse(true)

}
