package fileserver

import (
	"context"
	"encoding/json"
	"net/url"
)

type Statistics struct {
	ServerLoad     []float64 `json:"serverLoad"`
	DiskTotalSpace int64     `json:"diskTotalSpace"`
	DiskFreeSpace  int64     `json:"diskFreeSpace"`
}

func (s *Fileserver) GetStatistics(ctx context.Context) (*Statistics, error) {

	var query url.Values = map[string][]string{
		"action": {"statistics"},
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

	var stats *Statistics

	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		return nil, err
	}

	return stats, nil

}
