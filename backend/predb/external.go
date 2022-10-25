package predb

import (
	"atus/backend/category"
	"atus/backend/helpers"
	"atus/backend/request"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

// example: https://predb.ovh/api/v1/?q=%22Heat.1995.GERMAN.DL.2160P.UHD.BLURAY.X265-WATCHABLE%22
type Nuke struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type PreDBEntry struct {
	Name  string `json:"name"`
	Cat   string `json:"cat"`
	URL   string `json:"url"`
	PreAt int64  `json:"preAt"`
	Nuke  *Nuke  `json:"nuke"`
}

type PreDBResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		RowCount int           `json:"rowCount"`
		Rows     []*PreDBEntry `json:"rows"`
	} `json:"data"`
}

func getExternalData(ctx context.Context, rlsName string) (*PreDBEntry, error) {

	v := url.Values{}
	v.Set("q", `"`+url.QueryEscape(rlsName)+`"`)

	url := &url.URL{
		Scheme:   "https",
		Host:     "predb.ovh",
		Path:     "api/v1/",
		RawQuery: v.Encode(),
	}

	req, err := request.NewWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := req.Do()

	if err != nil {
		// sleep and retry if we reach the rate limit
		// if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(time.Second * 10)
		return getExternalData(ctx, rlsName)
		// }

		// return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pResp PreDBResp
	if err := json.Unmarshal(body, &pResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %s", err)
	}

	if pResp.Status != "success" {
		return nil, fmt.Errorf("predb returned error: %s", pResp.Message)
	}

	rlsNameComparable := helpers.ReplaceNonAlphanumeric(rlsName, "")
	for _, row := range pResp.Data.Rows {
		entryNameComparable := helpers.ReplaceNonAlphanumeric(row.Name, "")
		if entryNameComparable != "" && entryNameComparable == rlsNameComparable {
			return row, nil
		}
	}

	return nil, errors.New("no entries found")

}

var externalCategoryAssignments = map[category.Name][]string{
	category.XXX: {
		"xxx",
	},
	category.TV: {
		"tv", "serie",
	},
	category.Movie: {
		"movie", "screener", "vcd", "dvdr", "mvid", "divx", "xvid", "anime", "video",
	},
	category.Audio: {
		"flac", "mp3", "music", "audio", "soundtrack",
	},
	category.EBook: {
		"book",
	},
	category.App: {
		"app", "software", "0day", "0-day", "pda", "mac", "symbian", "iphone",
	},
	category.Docu: {
		"docu", "doku",
	},
	category.Game: {
		"console", "game", "psx", "ps1", "ps2", "ps3", "ps4", "ps5", "ps6", "ps7", "xbox", "dox", "nds", "gba", "ngc", "dreamcast", "docs", "gbc", "3ds",
	},
}

// predb.ovh is not using consistent category names, so we need to normalize them
func normalizeExternalCategory(categoryRaw, url string) category.Name {
	lowerCategory := strings.ToLower(categoryRaw)
	lowerURL := strings.ToLower(url)

	// predb.ovh - sometimes - returnes a url. we use the url to determine the category
	if lowerURL != "" {
		if strings.Contains(lowerURL, "soundcloud.com/") || strings.Contains(lowerURL, "discogs.com/") {
			return category.Audio
		}

		if strings.Contains(lowerURL, "tvmaze.com/shows/") || strings.Contains(lowerURL, "tvland.com/") || strings.Contains(lowerURL, "tvrage.com/") || strings.Contains(lowerURL, "thetvdb.com/") {
			return category.TV
		}
	}

	for c, aliases := range externalCategoryAssignments {
		for _, alias := range aliases {
			if strings.Contains(lowerCategory, alias) {
				return c
			}
		}
	}

	if lowerURL != "" {
		// category is still unknown but we found a imdb url. assume it's a movie
		if strings.Contains(lowerURL, "imdb.com/") {
			return category.Movie
		}
	}

	// if category contains "bluray" assume it's a movie
	// we might want to remove this in the future as a bluray could be a lot of different things
	if strings.Contains(lowerCategory, "bluray") {
		return category.Movie
	}

	return category.Unknown
}
