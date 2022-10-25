package fileserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type FileStatus struct {
	Index     int
	Completed bool
}

func (s *Fileserver) GetFileStatus(hash string, indices []int) ([]*FileStatus, error) {

	iStr := ""
	for _, i := range indices {
		if iStr != "" {
			iStr += ","
		}
		iStr += fmt.Sprintf("%d", i)
	}

	var query url.Values = map[string][]string{
		"action":  {"getFileStatus"},
		"hash":    {hash},
		"indices": {iStr},
	}

	ctx := context.Background()

	req, err := s.buildRequest(ctx, "GET", query, nil)
	if err != nil {
		return nil, err
	}

	resp, err := req.Do()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var r []*FileStatus
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil

}
