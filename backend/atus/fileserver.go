package atus

import (
	"atus/backend/config"
	"atus/backend/fileserver"
	"atus/backend/logger"
	"atus/backend/release"
	"atus/backend/scheduler"
	"context"
	"errors"
	"math/rand"
	"path"
	"sort"
	"sync"
	"time"
)

type Fileserver struct {
	updateStatisticsScheduler *scheduler.Scheduler
	updateListScheduler       *scheduler.Scheduler
	getMetaFilesTaskScheduler *scheduler.Scheduler
	listCache                 map[string]*fileserver.ListFile
	m                         sync.RWMutex

	*fileserver.Fileserver
}

// GetAllFileservers returns all fileservers
func (a *ATUS) GetAllFileservers() []*Fileserver {
	var servers []*Fileserver
	a.fileservers.Range(func(key, value interface{}) bool {
		servers = append(servers, value.(*Fileserver))
		return true
	})
	return servers
}

// AddNewFileserver adds a new fileserver
func (a *ATUS) AddNewFileserver(f *fileserver.Fileserver) error {
	err := f.Save()
	if err != nil {
		return err
	}

	fs := &Fileserver{
		Fileserver: f,
	}

	a.fileservers.Store(f.UID, fs)

	logger.Ref(logger.RefFileserver, f.UID).Type(logger.TypeFileserver).Infof("fileserver added")

	if err := fs.Enable(a); err != nil {
		return err
	}

	a.OnFileserversUpdated(fs)

	return nil
}

// GetFileserverByUID returns a fileserver by its UID
func (a *ATUS) GetFileserverByUID(uid string) *Fileserver {
	if fs, ok := a.fileservers.Load(uid); ok {
		return fs.(*Fileserver)
	}
	return nil
}

// DeleteFileserver deletes a fileserver
func (a *ATUS) DeleteFileserver(f *Fileserver) error {
	if err := f.Disable(a); err != nil {
		return err
	}

	a.fileservers.Delete(f.UID)

	logger.Ref(logger.RefFileserver, f.UID).Type(logger.TypeFileserver).Infof("fileserver deleted")

	a.OnFileserversUpdated(f)

	return f.Fileserver.Delete()
}

// Enable enables a fileserver
// Will spawn the necessary tasks and flags the fileserver as enabled
func (f *Fileserver) Enable(a *ATUS) error {

	f.updateStatisticsScheduler = scheduler.New(f.Fileserver.StatisticsInterval, func(ctx context.Context) {
		f.updateStatisticsTask(ctx, a)
	})
	f.updateStatisticsScheduler.Run(true)

	f.updateListScheduler = scheduler.New(f.Fileserver.ListInterval, func(ctx context.Context) {
		f.updateListTask(ctx, a)
	})
	f.updateListScheduler.Run(true)

	f.getMetaFilesTaskScheduler = scheduler.New(config.Base.Schedulers.FileserverGetMetaFilesInterval, func(ctx context.Context) {
		a.getMetaFilesTask(ctx, f)
	})
	f.getMetaFilesTaskScheduler.Run(true)

	f.Fileserver.Enabled = true

	logger.Ref(logger.RefFileserver, f.UID).Type(logger.TypeFileserver).Infof("fileserver enabled")

	a.OnFileserversUpdated(f)

	return f.Fileserver.Save()
}

// Disable disables a fileserver
func (f *Fileserver) Disable(a *ATUS) error {

	if f.updateStatisticsScheduler != nil {
		f.updateStatisticsScheduler.Stop()
	}

	if f.updateListScheduler != nil {
		f.updateListScheduler.Stop()
	}

	if f.getMetaFilesTaskScheduler != nil {
		f.getMetaFilesTaskScheduler.Stop()
	}

	f.Fileserver.Enabled = false

	logger.Ref(logger.RefFileserver, f.UID).Type(logger.TypeFileserver).Infof("fileserver disabled")

	a.OnFileserversUpdated(f)

	return f.Fileserver.Save()
}

type FileserverAllocationMethod string

const (
	FileserverAllocationMethodFill     = "FILL"
	FileserverAllocationMethodMostFree = "MOST_FREE"
	FileserverAllocationMethodRandom   = "RANDOM"
)

// GetFileserverForRelease selects a fileserver based on the allocation method
func (a *ATUS) GetFileserverForRelease(method FileserverAllocationMethod, release *Release) *Fileserver {

	servers := a.GetAllFileservers()

	// filter out disabled, unreachable, full or servers with insufficient space
	vi := 0
	for _, s := range servers {
		if !s.Fileserver.Enabled ||
			s.Fileserver.Statistics == nil ||
			s.Fileserver.Statistics.DiskFreeSpace < s.Fileserver.MinFreeDiskSpace ||
			s.Fileserver.Statistics.DiskFreeSpace < release.Size {
			continue
		}
		servers[vi] = s
		vi++
	}

	servers = servers[:vi]

	if len(servers) == 0 {
		return nil
	}

	if method == FileserverAllocationMethodRandom {
		rand.Seed(time.Now().UnixNano())
		return servers[rand.Intn(len(servers))]
	}

	// sort by allocation method
	sort.Slice(servers, func(i, j int) bool {
		if method == FileserverAllocationMethodFill {
			return servers[i].Fileserver.Statistics.DiskFreeSpace < servers[j].Fileserver.Statistics.DiskFreeSpace
		}

		return servers[i].Fileserver.Statistics.DiskFreeSpace > servers[j].Fileserver.Statistics.DiskFreeSpace
	})

	return servers[0]

}

// GetDownloadState returns the download state of a release
func (a *ATUS) GetDownloadState(fsUID, torrentHash string) (*fileserver.ListFile, error) {

	fs := a.GetFileserverByUID(fsUID)
	if fs == nil {
		return nil, errors.New("fileserver not found")
	}

	fs.m.RLock()
	defer fs.m.RUnlock()

	if _, ok := fs.listCache[torrentHash]; !ok {
		return nil, errors.New("torrent not found")
	}

	return fs.listCache[torrentHash], nil

}

// updateStatisticsTask updates the statistics of a fileserver (disk space, server load, etc.)
func (f *Fileserver) updateStatisticsTask(ctx context.Context, a *ATUS) {

	stats, err := f.Fileserver.GetStatistics(ctx)
	if err != nil {
		// ToDo: Mark fileserver as unreachable
		logger.Ref(logger.RefFileserver, f.UID).Type(logger.TypeFileserver).Errorf("failed to update statistics: %s", err)
		return
	}

	f.Fileserver.Statistics = stats

	a.OnFileserversUpdated(f)

}

func (f *Fileserver) updateListTask(ctx context.Context, a *ATUS) {

	list, err := f.Fileserver.GetList(ctx, config.GetString("FILESERVER__DOWNLOAD_LABEL"))
	if err != nil {
		// ToDo: Mark fileserver as unreachable
		logger.Ref(logger.RefFileserver, f.UID).Type(logger.TypeFileserver).Errorf("failed to update filelist: %s", err)
		return
	}

	for _, file := range list {
		old, ok := f.listCache[file.Hash]
		if !ok || (old.State != file.State || old.DownloadRate != file.DownloadRate || old.Done != file.Done || old.ETA != file.ETA) {
			a.OnDownloadStateChanged(file)
		}
	}

	newList := make(map[string]*fileserver.ListFile)
	for _, t := range list {
		newList[t.Hash] = t
	}

	f.listCache = newList

}

func (a *ATUS) getMetaFilesTask(ctx context.Context, f *Fileserver) {

	logWithRef := logger.Ref(logger.RefFileserver, f.UID).Type(logger.TypeFileserver)

	logWithRef.Debug("getMetaFilesTask started")

	a.pendingReleases.Range(func(k, v interface{}) bool {

		r := v.(*Release)

		// only get meta files for torrents that are currently - or just finished - downloading
		if r.State != release.StateDownloading && r.State != release.StateDownloaded {
			return true
		}

		// build a list of all files that need to be downloaded
		metaIndices := []int{}
		metaMap := make(map[int]*release.MetaFile)
		for _, mf := range r.MetaFiles {
			if mf.Index != -1 && mf.State == release.MetafileStateUnknown {
				metaIndices = append(metaIndices, mf.Index)
				metaMap[mf.Index] = mf
			}
		}

		anyUpdated := false
		if len(metaIndices) > 0 {

			logWithRef.Debugf("getFileStatus for %s (%v)", r.Hash, metaIndices)
			fileStatus, err := f.Fileserver.GetFileStatus(r.Hash, metaIndices)

			if err != nil {
				// ToDo: mark release as broken / fileserver as unreachable
				logger.Ref(logger.RefRelease, r.UID).Type(logger.TypeFileserver).Errorf("failed to get file status: %v", err)
				return true
			}

			logWithRef.Debugf("getFileStatus for %s (%v) returned %v", r.Hash, metaIndices, fileStatus)

			for _, fs := range fileStatus {

				if !fs.Completed {
					logWithRef.Debugf("file %d for %s is not completed", fs.Index, r.Hash)
					continue
				}
				mf, ok := metaMap[fs.Index]
				if !ok {
					logWithRef.Debugf("file %d for %s is not in map", fs.Index, r.Hash)
					continue
				}

				// We found a meta file that is completed, download it
				logWithRef.Debugf("file %d for %s is completed, downloading", fs.Index, r.Hash)
				err = f.Fileserver.DownloadFile(fs.Index, r.Hash, path.Join(config.Base.Folders.Data, mf.ReleaseUID, mf.FileName))
				if err != nil {
					// ToDo: mark release  / metafile as broken
					logWithRef.Errorf("failed to download file: %v", err)
					continue
				}

				// set new state
				newState := release.MetafileStateProcessed

				// samples need additional processing
				if mf.Type == release.MetafileTypeSampleVideo {
					newState = release.MetafileStateDownloaded

					logWithRef.Debugf("file %d for %s is a sample, adding to queue", fs.Index, r.Hash)

					// add to sample queue
					a.sampleQueue <- r
				}

				mf.State = newState

				if err := mf.Save(); err != nil {
					// ToDo: mark release  / metafile as broken
					logWithRef.Type(logger.TypeGeneric).Errorf("failed to save metafile: %v", err)
					continue
				}

				anyUpdated = true

				logger.Ref(logger.RefRelease, f.UID).Debugf("[FS %s] successfully downloaded metafile %s", f.UID, mf.FileName)
			}
		}

		if anyUpdated {
			a.OnMetaFilesUpdated(r)
		}

		return true

	})
}

func (a *ATUS) GetFileserverStatistics() []map[string]interface{} {

	var ret []map[string]interface{}

	a.fileservers.Range(func(k, v interface{}) bool {

		fs := v.(*Fileserver)

		ret = append(ret, map[string]interface{}{
			"uid":        fs.UID,
			"name":       fs.Name,
			"enabled":    fs.Enabled,
			"statistics": fs.Statistics,
		})

		return true

	})

	// sort by name or uid if name is the same
	sort.Slice(ret, func(i, j int) bool {
		if ret[i]["name"] == ret[j]["name"] {
			return ret[i]["uid"].(string) < ret[j]["uid"].(string)
		}
		return ret[i]["name"].(string) < ret[j]["name"].(string)
	})

	return ret

}
