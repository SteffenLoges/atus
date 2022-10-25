package fileserver

import (
	"atus/backend/helpers"
	"atus/backend/request"
	"atus/backend/sqlite"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"
)

type Fileserver struct {
	Enabled            bool
	Name               string
	URL                *url.URL
	UID                string
	ListInterval       time.Duration
	SumFilesDownloaded int64
	Statistics         *Statistics
	StatisticsInterval time.Duration
	MinFreeDiskSpace   int64
}

func New(u *url.URL) *Fileserver {
	return &Fileserver{
		UID:                sqlite.GenerateUID("fileservers"),
		URL:                u,
		ListInterval:       time.Second * 5,
		StatisticsInterval: time.Second * 10,
		MinFreeDiskSpace:   25 * helpers.GiB,
	}
}

func (s *Fileserver) buildRequest(ctx context.Context, method string, query url.Values, body io.Reader) (*request.Request, error) {

	u := *s.URL
	v := u.Query()
	for key, values := range query {
		for _, value := range values {
			v.Add(key, value)
		}
	}
	u.RawQuery = v.Encode()

	req, err := request.NewWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	return req, nil

}

// GetAll returns all sources from the database
// Do NOT call this function directly, use atus.GetAllFileservers() instead
func GetAll() ([]*Fileserver, error) {

	var servers []*Fileserver

	rows, err := sqlite.Conn.Query(
		`SELECT 
			uid,
			name,
			url,
			enabled,
			list_interval,
			statistics_interval,
			sum_files_downloaded,
			min_free_disk_space
		FROM fileservers
		ORDER BY name ASC`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		s := &Fileserver{}
		var urlRaw string

		err := rows.Scan(
			&s.UID,
			&s.Name,
			&urlRaw,
			&s.Enabled,
			&s.ListInterval,
			&s.StatisticsInterval,
			&s.SumFilesDownloaded,
			&s.MinFreeDiskSpace,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning sources: %s", err)
		}

		u, err := url.Parse(urlRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing url: %s", err)
		}

		s.URL = u

		servers = append(servers, s)
	}

	return servers, nil

}

// Deletes the server from the database
// Do NOT call this function directly, use atus.DeleteFileserver() instead
func (s *Fileserver) Delete() error {

	_, err := sqlite.Conn.Exec(
		`DELETE FROM fileservers WHERE uid = ?`,
		s.UID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Fileserver) Save() error {

	_, err := sqlite.Conn.Exec(
		`INSERT INTO fileservers
			(
				uid,
				name,
				url,
				list_interval,
				statistics_interval,
				sum_files_downloaded,
				min_free_disk_space
			) VALUES
			(?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(uid) DO UPDATE SET
			name = ?,
			url = ?,
			list_interval = ?,
			statistics_interval = ?,
			sum_files_downloaded = ?,
			min_free_disk_space = ?,
			enabled = ?`,
		s.UID,
		s.Name,
		s.URL.String(),
		s.ListInterval,
		s.StatisticsInterval,
		s.SumFilesDownloaded,
		s.MinFreeDiskSpace,
		s.Name,
		s.URL.String(),
		s.ListInterval,
		s.StatisticsInterval,
		s.SumFilesDownloaded,
		s.MinFreeDiskSpace,
		s.Enabled,
	)

	if err != nil {
		return err
	}

	return nil

}
