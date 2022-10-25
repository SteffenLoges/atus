package fileserver

import (
	"context"
	"encoding/json"
	"net/url"
)

type FileState string

const (
	FileStateStarted  FileState = "STARTED"
	FileStatePaused   FileState = "PAUSED"
	FileStateStopped  FileState = "STOPPED"
	FileStateHashing  FileState = "HASHING"
	FileStateChecking FileState = "CHECKING"
	FileStateError    FileState = "ERROR"
)

type ListFile struct {
	Hash         string    `json:"hash"`
	State        FileState `json:"state"`
	DownloadRate int64     `json:"downloadRate"`
	Done         float64   `json:"done"` // download status in percent
	ETA          int64     `json:"eta"`  // in seconds
}

func (s *Fileserver) GetList(ctx context.Context, label string) ([]*ListFile, error) {

	var query url.Values = map[string][]string{
		"action": {"list"},
		"label":  {label},
	}

	req, err := s.buildRequest(ctx, "GET", query, nil)
	if err != nil {
		return nil, err
	}

	resp, err := req.Do()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var list []*ListFile
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil

}
