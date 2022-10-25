package fileserver

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
)

// DownloadFile downloads a file from the fileserver
func (s *Fileserver) DownloadFile(index int, hash, savePath string) error {

	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	var query url.Values = map[string][]string{
		"action": {"downloadFile"},
		"hash":   {hash},
		"index":  {fmt.Sprintf("%d", index)},
	}

	req, err := s.buildRequest(context.Background(), "GET", query, nil)
	if err != nil {
		return err
	}

	resp, err := req.Do()
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	return err

}
