package atus

import (
	"atus/backend/logger"
	"atus/backend/release"
	"atus/backend/scheduler"
	"atus/backend/source"
	"context"
	"errors"
	"net/url"
	"sync"
	"time"
)

type Source struct {
	rssFeedScheduler *scheduler.Scheduler
	*source.Source
}

func (a *ATUS) GetAllSources() []*Source {

	var sources []*Source

	a.sources.Range(func(_, value interface{}) bool {
		sources = append(sources, value.(*Source))
		return true
	})

	return sources

}

func (a *ATUS) AddNewSource(s *source.Source) error {

	if err := s.Save(); err != nil {
		return err
	}

	source := &Source{
		Source: s,
	}

	a.sources.Store(s.UID, source)

	if err := source.Enable(a, true); err != nil {
		return err
	}

	return nil

}

func (a *ATUS) GetSourceByUID(uid string) *Source {

	if s, ok := a.sources.Load(uid); ok {
		return s.(*Source)
	}

	return nil

}

func (a *ATUS) DeleteSource(s *Source) error {

	if err := s.Disable(); err != nil {
		return err
	}

	a.sources.Delete(s.UID)

	return s.Source.Delete()

}

// simple cache to prevent unnecessary requests
var knownReleaseNames sync.Map

// checkRSSFeedTaskInit initializes the task that checks the RSS feed for new releases.
func (s *Source) checkRSSFeedTaskInit(releaseChan chan<- *release.Release) scheduler.Task {

	logWithRef := logger.Ref(logger.RefSource, s.UID).Type(logger.TypeSource)

	return func(ctx context.Context) {
		s.Source.LastCheck = time.Now()
		s.Source.NextCheck = time.Now().Add(s.Source.RSSInterval)
		s.Source.TimesChecked++
		defer s.Source.Save()

		// -- delete old entries from cache ---------
		// ToDo: This runs for every source. It should run only once per interval.
		// not a big deal, but still
		knownReleaseNames.Range(func(key, value interface{}) bool {
			if time.Since(value.(time.Time)) > 24*time.Hour {
				knownReleaseNames.Delete(key)
			}
			return true
		})

		// -- get feed ------------------------------
		feed, err := s.Source.GetRSSFeed(ctx)
		if err != nil {
			// ignore timeout errors
			if errors.Is(err, context.Canceled) {
				logWithRef.Debug("RSS feed check canceled (01)")
				return
			}

			// Other error
			// ToDo: Disable source and notify user
			// s.Enabled = false
			logWithRef.Errorf("Error getting RSS feed: %v", err)
			return
		}

		for _, item := range feed {

			logWithRef.Debugf("Checking item: %s", item.Title)

			// Check if the name found in the rss feed is already in the cache
			// so we don't need to download the meta file again.
			if _, ok := knownReleaseNames.Load(item.Title); ok {
				continue
			}

			knownReleaseNames.Store(item.Title, time.Now())

			// -- find urls ---------------------------
			metaURL, ok := s.Source.GetMetaURL(item)
			if !ok {
				logWithRef.Debugf("No meta url found in item: %s", item.Title)
				continue
			}

			var imageURL *url.URL
			if i, ok := s.Source.GetImageURL(item); ok {
				imageURL = i
			}

			// -- create release instance -------------
			release, err := release.New(ctx, s.Source, item.Title, metaURL, imageURL)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					logWithRef.Debug("RSS feed check canceled (02)")
					break
				}
				logWithRef.Errorf("Error creating release instance: %v", err)
				continue
			}

			// thats all we need for now, send the release to the release channel
			releaseChan <- release

		}
	}
}

func (s *Source) Enable(a *ATUS, runImmediate bool) error {

	task := s.checkRSSFeedTaskInit(a.releaseChan)
	s.rssFeedScheduler = scheduler.New(s.Source.RSSInterval, task)
	s.rssFeedScheduler.Run(runImmediate)

	s.Source.Enabled = true

	logger.Ref(logger.RefSource, s.UID).Type(logger.TypeSource).Info("Source enabled")

	return s.Source.Save()

}

func (s *Source) Disable() error {

	if s.rssFeedScheduler != nil {
		s.rssFeedScheduler.Stop()
	}

	s.Source.Enabled = false

	logger.Ref(logger.RefSource, s.UID).Type(logger.TypeSource).Info("Source disabled")

	return s.Source.Save()

}
