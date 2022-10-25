package source

import (
	"atus/backend/helpers"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/clbanning/mxj/v2"
	"golang.org/x/text/encoding/ianaindex"
)

type RSSFeedItem struct {
	ID       string
	Title    string
	MetaURL  string
	ImageURL string
}

func (s *Source) GetRSSFeed(ctx context.Context) ([]*ParsedFeedItem, error) {

	resp, err := s.MakeRequest(ctx, s.RSSURL.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	items, err := parseRSSFeed(resp.Body)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type ParsedFeedItem struct {
	ID    string
	Title string
	URLs  []*ParsedFeedItemURL
}

type ParsedFeedItemURL struct {
	Path string
	URL  *url.URL
}

func parseRSSFeed(r io.ReadCloser) ([]*ParsedFeedItem, error) {

	defer r.Close()

	mxj.CoerceKeysToLower(true)
	mxj.CustomDecoder = &xml.Decoder{
		// @see: https://dzhg.dev/posts/2020/08/how-to-parse-xml-with-non-utf8-encoding-in-go/
		CharsetReader: func(charset string, reader io.Reader) (io.Reader, error) {
			e, err := ianaindex.IANA.Encoding(charset)
			if err != nil {
				return nil, fmt.Errorf("encoding %s: %s", charset, err.Error())
			}
			if e == nil {
				// Assume it's compatible with (a subset of) UTF-8 encoding
				// Bug: https://github.com/golang/go/issues/19421
				return reader, nil
			}
			return e.NewDecoder().Reader(reader), nil
		},
	}

	reader, err := mxj.NewMapXmlReader(r)
	if err != nil {
		return nil, err
	}

	channelItems, err := reader.ValuesForPath("*.channel.item")
	if err != nil {
		return nil, err
	}

	var ParsedFeedItems []*ParsedFeedItem
	for _, ci := range channelItems {
		items, ok := ci.(map[string]interface{})
		if !ok {
			continue
		}

		ParsedFeedItem := ParsedFeedItem{}

		for key, item := range items {
			if key == "title" {
				ParsedFeedItem.Title = item.(string)
				continue
			}

			// some trackers use 'name' instead of 'title' so we'll check for both but prioritize 'title'
			if key == "name" && ParsedFeedItem.Title == "" {
				ParsedFeedItem.Title = item.(string)
				continue
			}

			// we do NOT rely on 'link' as it's not always present or the best one to use
			// some trackers use non-standard fields like 'ssl-link' or 'url'
			// we'll dump all URLs found in the item and check them later
			ParsedFeedItem.URLs = append(ParsedFeedItem.URLs, findURLsRecursive(key, item)...)
		}

		// ignore feed items without a link
		if len(ParsedFeedItem.URLs) > 0 {
			ParsedFeedItems = append(ParsedFeedItems, &ParsedFeedItem)
		}
	}

	// Find IDs
	for _, item := range ParsedFeedItems {
		id, err := item.findID()
		if err != nil {
			continue
		}
		item.ID = id
	}

	return ParsedFeedItems, nil
}

// findURLsRecursive recursively traverses the map and returns all valid URLs
func findURLsRecursive(key string, m interface{}) []*ParsedFeedItemURL {
	var findings []*ParsedFeedItemURL

	switch m := m.(type) {
	case map[string]interface{}:
		for k, v := range m {
			findings = append(findings, findURLsRecursive(key+"_"+k, v)...)
		}
	case string:
		if url, isValid := helpers.ValidateURL(m); isValid {
			findings = append(
				findings,
				&ParsedFeedItemURL{
					Path: key,
					URL:  url,
				},
			)
		}
	}

	return findings
}

// Matches the following formats:
//   - https://www.example.com/download.php?torrent=12345
//   - https://www.example.com/download.php?id=12345
//   - https://www.example.com/download/12345
//   - https://www.example.com/sldownload/12345
//   - https://www.example.com/download.php/12345/filename.torrent
var rssFeedIDRegExp = regexp.MustCompile(`download(?:(?:\.php)|-)?\??(?:(?:id|torrent)=|\/)?(?P<id>\d+)`)
var rssFeedIDRegExpIDIndex = rssFeedIDRegExp.SubexpIndex("id")

// last resort:
// use any number of digits as an ID.
// will result in a false positive if the url contains any parameters starting with digits before the actual ID
// e.g. https://www.example.com/download.php?ssl=1&torrent=12345
var rssFeedIDFallbackRegExp = regexp.MustCompile(`[\/|=](?P<id>\d+)`)
var rssFeedIDFallbackRegExpIDIndex = rssFeedIDFallbackRegExp.SubexpIndex("id")

func (p *ParsedFeedItem) findID() (string, error) {

	var fallbackID string
	for _, url := range p.URLs {

		// try to find an ID in the URL
		submatch := rssFeedIDRegExp.FindStringSubmatch(url.URL.String())
		if len(submatch) > 1 {
			return submatch[rssFeedIDRegExpIDIndex], nil
		}

		// try to find an ID in the URL using a fallback regexp
		submatch = rssFeedIDFallbackRegExp.FindStringSubmatch(url.URL.String())
		if len(submatch) > 1 {
			fallbackID = submatch[rssFeedIDFallbackRegExpIDIndex]
			// we don't return here because we want to keep checking for a better match
		}
	}

	if fallbackID != "" {
		return fallbackID, nil
	}

	return "", errors.New("no ID found")

}

func (p *ParsedFeedItem) GetURLFromGeneric(genericURL string) (*url.URL, bool) {
	parsed, err := url.Parse(genericURL)
	if err != nil {
		return nil, false
	}

	query := parsed.RawQuery
	query = strings.ReplaceAll(query, "{id}", p.ID)
	query = strings.ReplaceAll(query, "{title}", p.Title)
	ok := query != parsed.RawQuery
	parsed.RawQuery = query

	return parsed, ok
}

func (p *ParsedFeedItem) GetURLFromPath(path string) (*url.URL, bool) {
	for _, url := range p.URLs {
		if url.Path == path {
			return url.URL, true
		}
	}
	return nil, false
}
