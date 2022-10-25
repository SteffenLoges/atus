package source

import (
	"atus/backend/request"
	"atus/backend/sqlite"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Source struct {
	UID         string
	Name        string
	Favicon     string
	RSSURL      *url.URL
	RSSInterval time.Duration

	// time of last request to the source
	lastRequest time.Time

	// wait time between requests to the source
	RequestWaitTime time.Duration

	Cookies               []*http.Cookie
	MetaPath              string
	MetaPathUseAsKey      bool
	ImagePath             string
	ImagePathUseAsKey     bool
	LastCheck             time.Time
	NextCheck             time.Time
	TimesChecked          int64
	SumTorrentsDownloaded int64
	SumImagesDownloaded   int64
	SumReleasesDownloaded int64
	Enabled               bool
}

func New(u *url.URL, c []*http.Cookie) *Source {
	return &Source{
		UID:             sqlite.GenerateUID("sources"),
		RSSURL:          u,
		Cookies:         c,
		RequestWaitTime: time.Millisecond * 300,
		RSSInterval:     time.Minute * 2,
	}
}

// GetAll returns all sources from the database
// Do NOT call this function directly, use atus.GetAllSources() instead
func GetAll() ([]*Source, error) {

	var sources []*Source

	rows, err := sqlite.Conn.Query(
		`SELECT 
			uid,
			name,
			favicon,
			enabled,
			cookies,
			rss_url,
			rss_interval,
			request_waittime,
			last_check,
			meta_path,
			meta_path_use_as_key,
			image_path,
			image_path_use_as_key,
			times_checked,
			sum_torrents_downloaded,
			sum_images_downloaded,
			sum_releases_downloaded
		FROM sources
		ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var source Source
		var cookieBytes []byte
		var rssURLStr, lastCheckStr string

		err := rows.Scan(
			&source.UID,
			&source.Name,
			&source.Favicon,
			&source.Enabled,
			&cookieBytes,
			&rssURLStr,
			&source.RSSInterval,
			&source.RequestWaitTime,
			&lastCheckStr,
			&source.MetaPath,
			&source.MetaPathUseAsKey,
			&source.ImagePath,
			&source.ImagePathUseAsKey,
			&source.TimesChecked,
			&source.SumTorrentsDownloaded,
			&source.SumImagesDownloaded,
			&source.SumReleasesDownloaded,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning sources: %s", err)
		}

		if err := json.Unmarshal(cookieBytes, &source.Cookies); err != nil {
			return nil, fmt.Errorf("error unmarshalling cookies: %s", err)
		}

		source.RSSURL, err = url.Parse(rssURLStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing RSS URL for source %s: %s", source.UID, err)
		}

		if lastCheck, err := time.Parse(time.RFC3339, lastCheckStr); err == nil {
			source.LastCheck = lastCheck
		}

		sources = append(sources, &source)
	}

	return sources, nil

}

// Save saves the source to the database
// If the source already exists, it will be updated
func (s *Source) Save() error {
	rssURLStr := s.RSSURL.String()

	cooieStr, err := json.Marshal(s.Cookies)
	if err != nil {
		return err
	}

	_, err = sqlite.Conn.Exec(
		`INSERT INTO sources
			(
				uid, 
				name, 
				favicon, 
				cookies, 
				rss_url, 
				rss_interval, 
				request_waittime, 
				last_check, 
				meta_path, 
				meta_path_use_as_key, 
				image_path, 
				image_path_use_as_key,
				times_checked,
				sum_torrents_downloaded,
				sum_images_downloaded,
				sum_releases_downloaded
			) VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(uid) DO UPDATE SET
			name = ?,
			favicon = ?,
			enabled = ?,
			cookies = ?,
			rss_url = ?,
			rss_interval = ?,
			request_waittime = ?,
			last_check = ?,
			meta_path = ?,
			meta_path_use_as_key = ?,
			image_path = ?,
			image_path_use_as_key = ?,
			times_checked = ?,
			sum_torrents_downloaded = ?,
			sum_images_downloaded = ?,
			sum_releases_downloaded = ?`,
		s.UID,
		s.Name,
		s.Favicon,
		cooieStr,
		rssURLStr,
		s.RSSInterval,
		s.RequestWaitTime,
		s.LastCheck,
		s.MetaPath,
		s.MetaPathUseAsKey,
		s.ImagePath,
		s.ImagePathUseAsKey,
		s.TimesChecked,
		s.SumTorrentsDownloaded,
		s.SumImagesDownloaded,
		s.SumReleasesDownloaded,
		s.Name,
		s.Favicon,
		s.Enabled,
		cooieStr,
		rssURLStr,
		s.RSSInterval,
		s.RequestWaitTime,
		s.LastCheck.Format(time.RFC3339),
		s.MetaPath,
		s.MetaPathUseAsKey,
		s.ImagePath,
		s.ImagePathUseAsKey,
		s.TimesChecked,
		s.SumTorrentsDownloaded,
		s.SumImagesDownloaded,
		s.SumReleasesDownloaded,
	)

	if err != nil {
		return err
	}

	// delete favicon from tempFavicons to prevent it from being deleted
	for i, f := range tempFavicons {
		if f == s.Favicon {
			tempFavicons = append(tempFavicons[:i], tempFavicons[i+1:]...)
			break
		}
	}

	return nil

}

// Delete deletes the source from the database
// Do NOT call this function directly, use atus.Source.Delete() instead
func (s *Source) Delete() error {
	_, err := sqlite.Conn.Exec(
		`DELETE FROM sources WHERE uid = ?`,
		s.UID,
	)

	if err != nil {
		return err
	}

	return nil
}

// MakeRequest makes a request to the tracker with cookies
func (s *Source) MakeRequest(ctx context.Context, url string) (*http.Response, error) {

	// check if last request was too recent
	if time.Since(s.lastRequest) < s.RequestWaitTime {
		time.Sleep(s.RequestWaitTime - time.Since(s.lastRequest))
	}
	s.lastRequest = time.Now()

	req, err := request.NewWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range s.Cookies {
		req.Raw.AddCookie(cookie)
	}

	return req.Do()

}

func (s *Source) GetMetaURL(p *ParsedFeedItem) (*url.URL, bool) {
	if s.MetaPathUseAsKey {
		return p.GetURLFromPath(s.MetaPath)
	}

	return p.GetURLFromGeneric(s.MetaPath)
}

func (s *Source) GetImageURL(p *ParsedFeedItem) (*url.URL, bool) {
	if s.ImagePathUseAsKey {
		return p.GetURLFromPath(s.ImagePath)
	}

	return p.GetURLFromGeneric(s.ImagePath)
}
