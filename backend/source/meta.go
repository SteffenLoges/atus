package source

import (
	"atus/backend/logger"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/h2non/filetype"
)

// findImageFileURL tries to find an image file URL in the RSS feed item
// this will only be used during source-setup
// We migth at some point want to use this during crawling in case the tracker uses images from 3rd party sites
func (s *Source) FindImageFile(items []*ParsedFeedItem, maxRequests int) (*ParsedFeedItemURL, error) {

	// Keep track of the amount of requests we've made to the tracker
	blacklistedURLs := make(map[string]bool)
	requests := 0

	for _, item := range items {

		var urls []*ParsedFeedItemURL
		// filter and sort the urls by likelyhood of being a image file
		for _, url := range item.URLs {
			urlStr := url.URL.String()

			// ignore urls that have been already checked
			if blacklistedURLs[urlStr] {
				continue
			}

			// ignore urls that don't contain the ID
			// if !strings.Contains(urlStr, item.id) {
			// 	continue
			// }

			// ignore urls that end with .torrent
			if strings.HasSuffix(urlStr, ".torrent") {
				continue
			}

			// prioritize urls that end with .jpg, .png, .gif, .jpeg
			if strings.HasSuffix(urlStr, ".jpg") || strings.HasSuffix(urlStr, ".png") || strings.HasSuffix(urlStr, ".gif") || strings.HasSuffix(urlStr, ".jpeg") {
				urls = append([]*ParsedFeedItemURL{url}, urls...)
				continue
			}

			urls = append(urls, url)
		}

		for _, url := range urls {
			// netvision specific
			// prioritize full size images
			if strings.Contains(url.URL.String(), "f-"+item.ID) {
				urls = append([]*ParsedFeedItemURL{url}, urls...)
				continue
			}
		}

		for _, url := range urls {
			if requests >= maxRequests {
				return nil, errors.New("too many requests")
			}

			if requests > 0 {
				time.Sleep(time.Second)
			}

			requests++

			logger.Debugf("[NEW SOURCE] checking image url: %s", url.URL.String())

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if _, _, err := s.GetImage(ctx, url.URL.String()); err != nil {
				blacklistedURLs[url.URL.String()] = true
				continue
			}

			// we found an image file
			// return strings.ReplaceAll(url.URL, item.id, "{id}"), nil
			return url, nil

		}

	}

	return nil, errors.New("no image file url found")

}

// data, ext, err
func (s *Source) GetImage(ctx context.Context, url string) ([]byte, string, error) {

	resp, err := s.MakeRequest(context.Background(), url)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		return nil, "", fmt.Errorf("error getting filetype: %s", err)
	}

	if !filetype.IsImage(buf) {
		return nil, "", fmt.Errorf("not an image: %s", kind)
	}

	s.SumImagesDownloaded++

	return buf, "." + kind.Extension, nil

}

// FindTorrentFile tries to find a torrent file URL in the RSS feed item
func (s *Source) FindTorrentFile(items []*ParsedFeedItem, maxRequests int) (*ParsedFeedItemURL, error) {

	// Keep track of the amount of requests we've made to the tracker
	blacklistedURLs := make(map[string]bool)
	requests := 0

	for _, item := range items {

		var urls []*ParsedFeedItemURL
		// filter and sort the urls by likelyhood of being a meta file
		for _, url := range item.URLs {
			urlStr := url.URL.String()

			// ignore urls that have been already checked
			if blacklistedURLs[urlStr] {
				continue
			}

			// ignore urls that don't contain the ID
			// if !strings.Contains(url, item.id) {
			// 	continue
			// }

			// url.to/{id}/{name}.torrent
			if strings.HasSuffix(urlStr, ".torrent") {
				urls = append([]*ParsedFeedItemURL{url}, urls...)
				continue
			}

			// url.to/download.php?torrent=....
			if strings.Contains(urlStr, "/download") {
				urls = append([]*ParsedFeedItemURL{url}, urls...)
				continue
			}

			urls = append(urls, url)
		}

		for _, url := range urls {
			if requests >= maxRequests {
				return nil, errors.New("too many requests")
			}

			if requests > 0 {
				time.Sleep(time.Second)
			}

			requests++

			logger.Debugf("[NEW SOURCE] checking torrent url: %s", url.URL.String())

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			_, err := s.GetTorrentFile(ctx, url.URL.String())
			if err != nil {
				blacklistedURLs[url.URL.String()] = true
				continue
			}

			// we found a meta file
			// return strings.ReplaceAll(url.URL, item.id, "{id}"), nil
			return url, nil
		}

	}

	return nil, errors.New("no torrent url found")

}

func (s *Source) GetTorrentFile(ctx context.Context, url string) ([]byte, error) {

	resp, err := s.MakeRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")

	if !strings.Contains(contentType, "application/x-bittorrent") {
		return nil, fmt.Errorf("Content-Type is %s, expected application/x-bittorrent\n\nresponse: %s", contentType, buf)
	}

	s.SumTorrentsDownloaded++

	return buf, nil

}
