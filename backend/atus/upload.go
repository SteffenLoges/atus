package atus

import (
	"atus/backend/bencode"
	"atus/backend/config"
	"atus/backend/release"
	"atus/backend/request"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

// returns (newTorrentFile, error)
func (a *ATUS) UploadReleaseToTracker(ctx context.Context, r *Release) (*bencode.Dict, error) {

	// get torrent & nfo file
	var torrent, nfo []byte
	for _, mf := range r.MetaFiles {
		if mf.Type != release.MetafileTypeTorrent && mf.Type != release.MetafileTypeNFO {
			continue
		}

		file, err := mf.GetFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get %s file for release %s: %s", strings.ToLower(string(mf.Type)), r.Name, err.Error())
		}

		if mf.Type == release.MetafileTypeTorrent {
			torrent = file
			continue
		}

		if mf.Type == release.MetafileTypeNFO {
			nfo = file
			continue
		}
	}

	d := &Destination{
		TrackerAnnounceURL: config.GetString("UPLOAD__TRACKER_ANNOUNCE_URL"),
		Comment:            config.GetString("UPLOAD__COMMENT"),
		CreatedBy:          config.GetString("UPLOAD__CREATED_BY"),
		UserID:             config.GetString("UPLOAD__USER_ID"),
		APIURL:             config.GetString("UPLOAD__API_URL"),
		APIAuthToken:       config.GetString("API__AUTH_TOKEN"),
	}

	return d.UploadRelease(ctx, r, torrent, nfo)

}

type Destination struct {
	TrackerAnnounceURL string
	Comment            string
	CreatedBy          string
	UserID             string
	APIURL             string
	APIAuthToken       string
}

// uploadRelease uploads a release to the destination tracker
func (d *Destination) UploadRelease(ctx context.Context, r *Release, torrent, nfo []byte) (*bencode.Dict, error) {

	// build torrent file for tracker
	dict, err := bencode.BDecode(torrent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode torrent file for release %s: %s", r.Name, err.Error())
	}

	// set values
	dict.Announce = d.TrackerAnnounceURL
	dict.Comment = d.Comment
	dict.Private = 1
	dict.CreatedBy = d.CreatedBy
	dict.CreatedAt = time.Now().Unix()
	dict.Bot = "ATUS (github.com/SteffenLoges/atus)"

	encodedTorrent, err := dict.BEncode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode torrent file for release %s: %s", r.Name, err.Error())
	}

	// hash
	hash, err := dict.GenHash()
	if err != nil {
		return nil, err
	}

	// build request
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// add files
	files := map[string][]byte{
		"torrent": encodedTorrent,
		"nfo":     nfo,
	}

	for name, data := range files {
		w, err := writer.CreateFormFile(name, fmt.Sprintf("%s.%[1]s", name))
		if err != nil {
			return nil, err
		}

		if _, err := w.Write(data); err != nil {
			return nil, err
		}
	}

	// -- post data ---------------------------------------------------------------------------------
	postData := map[string]string{}

	// fileList
	fileList := dict.GetFiles()
	sortedFileList := bencode.SortFiles(fileList)

	marshaledFileList, err := json.Marshal(sortedFileList)
	if err != nil {
		return nil, err
	}

	postData["fileList"] = string(marshaledFileList)

	// metaFiles
	marshaledMetaFiles, err := json.Marshal(r.MetaFiles)
	if err != nil {
		return nil, err
	}

	postData["metaFiles"] = string(marshaledMetaFiles)
	postData["hash"] = hash
	postData["name"] = r.Name
	postData["category"] = r.Category
	postData["categoryRaw"] = r.CategoryRaw
	postData["pre"] = r.Pre.UTC().String()
	postData["userID"] = d.UserID

	// -- send request ------------------------------------------------------------------------------
	for k, v := range postData {
		if err := writer.WriteField(k, v); err != nil {
			panic(err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s?action=upload&authentication=%s", d.APIURL, d.APIAuthToken)
	req, err := request.NewWithContext(ctx, "POST", url, buf)
	if err != nil {
		return nil, err
	}

	req.Raw.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := req.Do()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	//  read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	var respStruct struct {
		Success bool
		Message string
	}

	if err := json.Unmarshal(body, &respStruct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s; err: %s", body, err.Error())
	}

	// check if error is set
	if !respStruct.Success {
		return nil, fmt.Errorf("failed to upload release: %s; Raw: %s", respStruct.Message, body)
	}

	return dict, nil

}
