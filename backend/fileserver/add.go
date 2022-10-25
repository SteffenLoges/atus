package fileserver

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/url"
)

type AddTorrentResponse struct {
	Hash        string
	Prioritized []interface{}
	Debug       interface{}
}

// AddTorrent adds a torrent to the fileserver
func (s *Fileserver) AddTorrent(ctx context.Context, file []byte, name, label string) (*AddTorrentResponse, error) {

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	w, err := writer.CreateFormFile("meta", name)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(file); err != nil {
		return nil, err
	}

	if err := writer.WriteField("label", label); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	var query url.Values = map[string][]string{
		"action": {"add"},
	}

	req, err := s.buildRequest(ctx, "POST", query, buf)
	if err != nil {
		return nil, err
	}

	req.Raw.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := req.Do()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var r *AddTorrentResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil

}
