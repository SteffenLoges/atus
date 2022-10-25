package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/helpers"
	"atus/backend/logger"
	"atus/backend/source"
	"atus/backend/websocket"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var settingsSourcesAddCache sync.Map

type settingsSourcesAdd struct {
	source    *source.Source
	feedItems []*source.ParsedFeedItem
}

func Settings__SourcesAdd_SetRSSURL(r *websocket.Request) {

	var req struct {
		URL     string
		Cookies []*http.Cookie
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(fmt.Errorf("rss url is not a valid url: %s", err.Error()))
		return
	}

	logger.Debugf("[NEW SOURCE] rss url: %s", parsedURL.String())

	s := source.New(parsedURL, req.Cookies)

	// -- reading rss data --------------------------------------------------------------------------

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	feedItems, err := s.GetRSSFeed(ctx)
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(fmt.Sprintf("error while reading rss data: %s", err.Error()))
		return
	}

	logger.Debugf("[NEW SOURCE] found %d feed items", len(feedItems))

	// This ist most likely because the rss url is not valid
	if len(feedItems) == 0 {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse("invalid or empty rss feed")
		return
	}

	// -- at this point we have a valid rss feed, time to get some meta data ------------------------

	// -- get name --------------------------------
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	metaData, err := s.GetWebsiteMetadata(ctx)
	if err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(fmt.Errorf("couldn't get metadata from source: %s", err.Error()))
		return
	}

	logger.Debugf("[NEW SOURCE] found meta data: %+v", metaData)

	s.Name = metaData.Name

	// -- download favicon ------------------------
	if metaData.Favicon != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		faviconName, err := s.DownloadFavicon(ctx, metaData.Favicon)

		// We don't care about errors here, we just don't have a favicon
		if err == nil {
			s.Favicon = faviconName
			logger.Debugf("[NEW SOURCE] downloaded favicon: %s", faviconName)
		}
	}

	// ----------------------------------------------------------------------------------------------

	sas := settingsSourcesAdd{
		source:    s,
		feedItems: feedItems,
	}

	// max. number of requests to the tracker to find a meta or image file
	const maxRequests = 5

	metaPathAutoDetected := false
	if torrentFile, err := s.FindTorrentFile(feedItems, maxRequests); err == nil {
		sas.source.MetaPath = torrentFile.Path
		sas.source.MetaPathUseAsKey = true
		metaPathAutoDetected = true
		logger.Debugf("[NEW SOURCE] found meta file: %s", torrentFile.Path)
	} else {
		logger.Debugf("[NEW SOURCE] couldn't find meta file: %s", err.Error())
	}

	imagePathAutoDetected := false
	if imageFile, err := s.FindImageFile(feedItems, maxRequests); err == nil {
		sas.source.ImagePath = imageFile.Path
		sas.source.ImagePathUseAsKey = true
		imagePathAutoDetected = true
		logger.Debugf("[NEW SOURCE] found image file: %s", imageFile.Path)
	} else {
		logger.Debugf("[NEW SOURCE] couldn't find image file: %s", err.Error())
	}

	settingsSourcesAddCache.Store(sas.source.UID, sas)

	// ----------------------------------------------------------------------------------------------

	r.MarshalAndSendResponse(map[string]interface{}{
		"uid":                   sas.source.UID,
		"name":                  sas.source.Name,
		"favicon":               sas.source.Favicon,
		"metaPath":              sas.source.MetaPath,
		"metaPathUseAsKey":      sas.source.MetaPathUseAsKey,
		"metaPathAutoDetected":  metaPathAutoDetected,
		"imagePath":             sas.source.ImagePath,
		"imagePathUseAsKey":     sas.source.ImagePathUseAsKey,
		"imagePathAutoDetected": imagePathAutoDetected,
	})

}

func Settings__SourcesAdd_SetMetaPath(r *websocket.Request) {

	var req struct {
		UID              string
		MetaPath         string
		MetaPathUseAsKey bool
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	bi, ok := settingsSourcesAddCache.Load(req.UID)
	if !ok {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse("source not found, please start over")
		return
	}
	sas := bi.(settingsSourcesAdd)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	if sas.source.MetaPath == "" || sas.source.MetaPath != req.MetaPath {
		sas.source.MetaPath = req.MetaPath
		sas.source.MetaPathUseAsKey = req.MetaPathUseAsKey
		metaURL := ""
		for _, item := range sas.feedItems {
			if u, ok := sas.source.GetMetaURL(item); ok {
				metaURL = u.String()
				break
			}
		}

		if metaURL == "" {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse("we couldn't find a meta url with the given parameters")
			return
		}

		if _, isValid := helpers.ValidateURL(metaURL); !isValid {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse("the found meta url is not a valid url")
			return
		}

		_, err := sas.source.GetTorrentFile(ctx, metaURL)
		if err != nil {
			r.SetResponseCode(http.StatusInternalServerError)
			r.MarshalAndSendResponse(fmt.Sprintf("we couldn't get a valid meta file. Error: %s", err.Error()))
			return
		}

		sas.source.MetaPath = req.MetaPath
		sas.source.MetaPathUseAsKey = req.MetaPathUseAsKey
	}

	r.MarshalAndSendResponse(true)

}

func Settings__SourcesAdd_SetImagePath(r *websocket.Request) {

	var req struct {
		UID               string
		ImagePath         string
		ImagePathUseAsKey bool
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	bi, ok := settingsSourcesAddCache.Load(req.UID)
	if !ok {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse("source not found, please start over")
		return
	}
	sas := bi.(settingsSourcesAdd)

	// we cannot validate the image path here, because it might be possible that none of
	// current the feed items has an image

	sas.source.ImagePath = req.ImagePath
	sas.source.ImagePathUseAsKey = req.ImagePathUseAsKey

	r.MarshalAndSendResponse(true)

}

func Settings__SourcesAdd_SetSettings(r *websocket.Request) {

	var req struct {
		UID         string
		Name        string
		RSSInterval int64
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	bi, ok := settingsSourcesAddCache.Load(req.UID)
	if !ok {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse("source not found, please start over")
		return
	}
	sas := bi.(settingsSourcesAdd)

	sas.source.Name = req.Name
	sas.source.RSSInterval = time.Duration(req.RSSInterval) * time.Second

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	if err := a.AddNewSource(sas.source); err != nil {
		r.SetResponseCode(http.StatusInternalServerError)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	atus.SetSetupStepDone(r.Hub, atus.SetupStepSourceAdded)

	settingsSourcesAddCache.Delete(req.UID)

	r.MarshalAndSendResponse(true)

}
