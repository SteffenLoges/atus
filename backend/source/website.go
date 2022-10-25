package source

import (
	"atus/backend/config"
	"atus/backend/request"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/h2non/filetype"
)

type WebsiteMetadata struct {
	Name    string
	Favicon *url.URL
}

func (s *Source) GetWebsiteMetadata(ctx context.Context) (*WebsiteMetadata, error) {

	req, err := request.NewWithContext(ctx, "GET", s.RSSURL.Scheme+"://"+s.RSSURL.Host+"/", nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range s.Cookies {
		req.Raw.AddCookie(cookie)
	}

	resp, err := req.Do()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Find the <title> tag
	title := doc.Find("title").Text()
	if title == "" {
		// No title found, use host as name
		title = resp.Request.URL.Host
	}

	// find favicon
	icon := doc.Find("link[rel*=icon]").AttrOr("href", "")

	if icon == "" {
		// No icon meta tag found, use default
		icon = "/favicon.ico"
	}

	favicon, err := url.Parse(icon)
	if err != nil {
		return nil, err
	}

	// fix relative links
	if favicon.Scheme == "" {
		favicon.Scheme = resp.Request.URL.Scheme
	}

	if favicon.Host == "" {
		favicon.Host = resp.Request.URL.Host
	}

	return &WebsiteMetadata{
		Name:    title,
		Favicon: favicon,
	}, nil

}

// Store path to favicons in memory
// entry will be deleted after source is added to database
// remaining favicons will be deleted from disk on atus shutdown
var tempFavicons []string

// returns (the new filename, error)
func (s *Source) DownloadFavicon(ctx context.Context, u *url.URL) (string, error) {

	req, err := request.NewWithContext(ctx, "GET", u.String(), nil)

	if err != nil {
		return "", err
	}

	resp, err := req.Do()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		return "", err
	}

	if !filetype.IsImage(buf) {
		return "", fmt.Errorf("invalid file type: %s", kind)
	}

	file := s.UID + "." + kind.Extension

	// create favicons folder if it doesn't exist
	if _, err := os.Stat(path.Join(config.Base.Folders.Data, "favicons")); os.IsNotExist(err) {
		os.Mkdir(path.Join(config.Base.Folders.Data, "favicons"), 0755)
	}

	out, err := os.Create(path.Join(config.Base.Folders.Data, "favicons", file))
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = out.Write(buf)
	if err != nil {
		return "", err
	}

	tempFavicons = append(tempFavicons, file)

	return file, nil

}

func DeleteTempFavicons() {
	for _, file := range tempFavicons {
		os.Remove(path.Join(config.Base.Folders.Data, "favicons", file))
	}
}
