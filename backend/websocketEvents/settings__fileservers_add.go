package websocketEvents

import (
	"atus/backend/atus"
	"atus/backend/fileserver"
	"atus/backend/helpers"
	"atus/backend/websocket"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"
)

var settingsFileserversAddCache sync.Map

func Settings__FileserversAdd_SetURL(r *websocket.Request) {

	var req struct {
		URL string
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(fmt.Sprintf("invalid URL: %s", err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fs := fileserver.New(parsedURL)
	fs.Name = parsedURL.Hostname()
	stats, err := fs.GetStatistics(ctx)
	if err != nil {
		orgErr := err

		// request was unsuccessful, add default plugin paths to the URL and try again
		for _, pluginPath := range []string{"plugins/atus/api.php"} {
			time.Sleep(time.Millisecond * 300)
			newURL := parsedURL
			newURL.Path = path.Join(parsedURL.Path, pluginPath)
			fs.URL = newURL

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			stats, err = fs.GetStatistics(ctx)
			if err == nil {
				break
			}
		}

		// if still unsuccessful, return the original error
		if err != nil {
			r.SetResponseCode(http.StatusBadRequest)
			r.MarshalAndSendResponse(fmt.Sprintf("could not connect to fileserver. Make sure the server is running and the atus plugin is installed. err: %s", orgErr.Error()))
			return
		}
	}

	fs.Statistics = stats

	settingsFileserversAddCache.Store(fs.UID, fs)

	r.MarshalAndSendResponse(map[string]interface{}{
		"uid":                fs.UID,
		"name":               fs.Name,
		"listInterval":       fs.ListInterval / time.Second,
		"minFreeDiskSpace":   fs.MinFreeDiskSpace / helpers.GiB,
		"statisticsInterval": fs.StatisticsInterval / time.Second,
		"diskFreeSpace":      stats.DiskFreeSpace,
		"diskTotalSpace":     stats.DiskTotalSpace,
		"serverLoad":         stats.ServerLoad,
	})
}

func Settings__FileserversAdd_SetSettings(r *websocket.Request) {

	a := r.Hub.Ctx.Value(atus.ContextKey).(*atus.ATUS)

	var req struct {
		UID                string
		Name               string
		ListInterval       int64
		StatisticsInterval int64
		MinFreeDiskSpace   int64
	}

	if err := json.Unmarshal(r.Payload, &req); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(err.Error())
		return
	}

	sfc, ok := settingsFileserversAddCache.Load(req.UID)
	if !ok {
		r.SetResponseCode(http.StatusNotFound)
		r.MarshalAndSendResponse("fileserver not found, please start over")
		return
	}
	fs := sfc.(*fileserver.Fileserver)

	if req.Name != "" {
		fs.Name = req.Name
	}

	listInterval := time.Duration(req.ListInterval) * time.Second
	if listInterval > time.Second {
		fs.ListInterval = listInterval
	}

	statisticsInterval := time.Duration(req.StatisticsInterval) * time.Second
	if statisticsInterval > time.Second {
		fs.StatisticsInterval = statisticsInterval
	}

	if req.MinFreeDiskSpace > 0 {
		fs.MinFreeDiskSpace = helpers.GiB * req.MinFreeDiskSpace
	}

	if err := a.AddNewFileserver(fs); err != nil {
		r.SetResponseCode(http.StatusBadRequest)
		r.MarshalAndSendResponse(fmt.Sprintf("could not save settings: %s", err.Error()))
		return
	}

	atus.SetSetupStepDone(r.Hub, atus.SetupStepFileserverAdded)

	settingsFileserversAddCache.Delete(req.UID)

	r.MarshalAndSendResponse(true)
}
