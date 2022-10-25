package atus

import (
	"atus/backend/category"
	"atus/backend/config"
	"atus/backend/fileserver"
	"atus/backend/helpers"
	"atus/backend/logger"
	"atus/backend/release"
	"atus/backend/scheduler"
	"atus/backend/source"
	"atus/backend/video"
	"fmt"
	"sync"
	"time"
)

var ContextKey helpers.ContextKey = "atus"

type ATUS struct {
	pendingReleases sync.Map
	sources         sync.Map
	categories      sync.Map
	fileservers     sync.Map
	sampleQueue     chan *Release
	releaseChan     chan *release.Release

	OnReleaseAdded         func(*release.Release)
	OnReleaseStateUpdated  func(*Release, time.Time)
	OnMetaFilesUpdated     func(*Release)
	OnFileserversUpdated   func(*Fileserver)
	OnDownloadStateChanged func(*fileserver.ListFile)
}

func New() (*ATUS, error) {

	// -- check dependencies (ffmpeg, ffprobe, mp4dash)
	if err := video.CheckDependencies(); err != nil {

		fmt.Println(err.Error())
		time.Sleep(time.Hour)

		return nil, err
	}

	// -- init instance and create channels -------
	a := &ATUS{
		releaseChan: make(chan *release.Release, 500),

		// a channel size of 100 is to large for a sample queue.
		// the server clearly can't handle the load if there are that many samples in the queue
		// ToDo: warn the user if there are more then x pending samples
		sampleQueue: make(chan *Release, 100),

		OnReleaseAdded:         func(*release.Release) {},
		OnReleaseStateUpdated:  func(*Release, time.Time) {},
		OnMetaFilesUpdated:     func(*Release) {},
		OnFileserversUpdated:   func(*Fileserver) {},
		OnDownloadStateChanged: func(*fileserver.ListFile) {},
	}

	// -- get categories --------------------------
	allCategories, err := category.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not load categories: %s", err)
	}

	for _, c := range allCategories {
		a.categories.Store(c.Name, c)
	}

	// -- get pending releases --------------------
	pendingReleases, err := a.loadPendingReleases()
	if err != nil {
		return nil, fmt.Errorf("could not load pending releases: %s", err)
	}

	for _, r := range pendingReleases {
		a.pendingReleases.Store(r.Hash, r)
	}

	// -- get fileservers -------------------------
	allFileservers, err := fileserver.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not load fileservers: %s", err)
	}

	for _, f := range allFileservers {
		s := &Fileserver{
			Fileserver: f,
			listCache:  make(map[string]*fileserver.ListFile),
		}

		a.fileservers.Store(f.UID, s)

		// start scheduler if fileserver is enabled
		if s.Enabled {
			if err := s.Enable(a); err != nil {
				return nil, fmt.Errorf("could not enable fileserver: %s", err)
			}
		}
	}

	// -- get sources -----------------------------
	allSources, err := source.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not load sources: %s", err)
	}

	for _, s := range allSources {
		rs := &Source{
			Source: s,
		}
		a.sources.Store(s.UID, rs)

		// start scheduler if source is enabled
		if s.Enabled {
			if err := rs.Enable(a, false); err != nil {
				return nil, fmt.Errorf("could not enable source: %s", err)
			}
		}
	}

	// -- watch for new releases --------------------------------------------------------------------

	// unchecked releases channel
	go func() {
		for r := range a.releaseChan {
			a.onNewRelease(r)
		}
	}()

	// sample queue channel
	go func() {
		for r := range a.sampleQueue {
			if err := a.onNewSample(r); err != nil {
				logger.Ref(logger.RefRelease, r.UID).Type(logger.TypeSample).Errorf("error saving sample meta file: %s", err)
			}
		}
	}()

	// -- start pending release scheduler -----------------------------------------------------------
	pendingReleaseScheduler := scheduler.New(config.Base.Schedulers.ProcessPendingReleasesInterval, a.processPendingReleasesTask)
	pendingReleaseScheduler.Run(false)

	return a, nil

}
