package atus

import (
	"atus/backend/config"
	"atus/backend/logger"
	"atus/backend/predb"
	"atus/backend/release"
	"atus/backend/sqlite"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ToDo: We currently have release.Release and atus.Release.
// This is confusing and should be refactored.
type Release struct {
	UID           string
	Name          string
	Category      string
	CategoryRaw   string
	Hash          string
	Size          int64
	Pre           time.Time
	State         release.ReleaseState
	SourceUID     string
	MetaFiles     []*release.MetaFile
	FileserverUID string
}

func (a *ATUS) loadPendingReleases() ([]*Release, error) {

	rows, err := sqlite.Conn.Query(
		`SELECT 
			uid,
			name,
			hash,
			size,
			state,
			category, 
			category_raw,
			source_uid,
			fileserver_uid,
			pre
		FROM 
			releases 
		WHERE 
			state NOT IN(?, ?, ?)`,
		release.StateUploaded,
		release.StateGeneralError,
		release.StateUploadError,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pendingReleases []*Release
	for rows.Next() {
		r := &Release{}
		var preStr string
		if err := rows.Scan(
			&r.UID,
			&r.Name,
			&r.Hash,
			&r.Size,
			&r.State,
			&r.Category,
			&r.CategoryRaw,
			&r.SourceUID,
			&r.FileserverUID,
			&preStr,
		); err != nil {
			return nil, err
		}

		if lastCheck, err := time.Parse(time.RFC3339, preStr); err == nil {
			r.Pre = lastCheck
		}

		pendingReleases = append(pendingReleases, r)
	}

	// meta files
	for _, pr := range pendingReleases {
		upm, err := release.GetMetaFiles(pr.UID, release.MetaFileState(""))
		if err != nil {
			return nil, err
		}

		pr.MetaFiles = upm

		// put unprocessed samples in queue
		for _, m := range upm {
			if m.Type == release.MetafileTypeSampleVideo && m.State == release.MetafileStateDownloaded {
				a.sampleQueue <- pr
			}
		}
	}

	return pendingReleases, nil

}

func (a *ATUS) GetReleaseByUID(uid string) *Release {

	r := &Release{}
	var preStr string
	if err := sqlite.Conn.QueryRow(
		`SELECT 
			uid,
			name,
			hash,
			size,
			state,
			category, 
			category_raw,
			source_uid,
			fileserver_uid,
			pre
		FROM 
			releases 
		WHERE 
			uid = ?`,
		uid,
	).Scan(
		&r.UID,
		&r.Name,
		&r.Hash,
		&r.Size,
		&r.State,
		&r.Category,
		&r.CategoryRaw,
		&r.SourceUID,
		&r.FileserverUID,
		&preStr,
	); err != nil {
		return nil
	}

	if lastCheck, err := time.Parse(time.RFC3339, preStr); err == nil {
		r.Pre = lastCheck
	}

	// meta files
	upm, err := release.GetMetaFiles(r.UID, release.MetafileStateProcessed)
	if err != nil {
		return nil
	}

	r.MetaFiles = upm

	return r
}

func (a *ATUS) GetPendingReleaseByHash(hash string) *Release {
	if hash, ok := a.pendingReleases.Load(hash); ok {
		return hash.(*Release)
	}

	return nil
}

func (a *ATUS) DeleteRelease(uid, hash string) error {

	// Delete from pending releases, if exists
	if pr := a.GetPendingReleaseByUID(uid); pr != nil {
		a.pendingReleases.Delete(pr.Hash)
	}

	// Delete from database
	_, err := sqlite.Conn.Exec("DELETE FROM releases WHERE uid = ?", uid)
	if err != nil {
		return err
	}

	// delete meta files from database
	_, err = sqlite.Conn.Exec("DELETE FROM release_metafiles WHERE release_uid = ?", uid)
	if err != nil {
		return err
	}

	// Delete data folder
	os.RemoveAll(filepath.Join(config.Base.Folders.Data, uid))

	return nil
}

func (a *ATUS) GetPendingReleaseByUID(uid string) *Release {
	var pr *Release
	a.pendingReleases.Range(func(key, value interface{}) bool {
		if value.(*Release).UID == uid {
			pr = value.(*Release)
			return false
		}

		return true
	})

	return pr
}

func (a *ATUS) updatePendingReleaseState(r *Release, state release.ReleaseState) error {

	if r.State != state {
		logger.Ref(logger.RefRelease, r.UID).Infof("release state changed to %s", state)
	}

	r.State = state

	// remove release from pending releases if state is uploaded or error
	if state == release.StateUploaded || strings.HasSuffix(string(state), "ERROR") {
		a.pendingReleases.Delete(r.UID)
	}

	query := `UPDATE releases SET state = ?`
	binds := []interface{}{state}

	var uploadDate time.Time
	if state == release.StateUploaded {
		uploadDate = time.Now()
		query += `, uploaded = ?`
		binds = append(binds, uploadDate.Format(time.RFC3339))
	}

	_, err := sqlite.Conn.Exec(query+` WHERE uid = ?`, append(binds, r.UID)...)

	if err != nil {
		return err
	}

	a.OnReleaseStateUpdated(r, uploadDate)

	return nil
}

func (a *ATUS) assignFileserver(r *Release, fs *Fileserver) error {
	r.FileserverUID = fs.UID

	_, err := sqlite.Conn.Exec(`UPDATE releases SET fileserver_uid = ? WHERE uid = ?`, fs.UID, r.UID)
	return err
}

func (a *ATUS) onNewRelease(r *release.Release) {

	logWithRef := logger.Ref(logger.RefRelease, r.UID).Type(logger.TypeRelease)

	logWithRef.Infof("found new release %s on %s", r.Name, r.Source.Name)

	// -- Check if release is already in database -
	// ToDo: Allow to add file anyway to increase download speed
	if r.IsKnown() {
		logWithRef.Debugf("release is already in database")
		return
	}

	// -- Check if release is in predb ------------
	logWithRef.Type(logger.TypePredb).Debugf("checking predb for release %s", r.Name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pre, err := predb.GetRelease(ctx, r.Name)
	if err != nil {
		logWithRef.Type(logger.TypePredb).Errorf("error while getting release from predb: %s", err.Error())
		return
	}
	logWithRef.Type(logger.TypePredb).Debugf("release %s found in predb. PreTime: %s", r.Name, pre.At.String())

	// -- Check age -------------------------------
	if config.GetInt64("FILTERS__MAX_AGE") > 0 {
		if time.Since(pre.At) > time.Duration(config.GetInt64("FILTERS__MAX_AGE"))*time.Minute {
			logWithRef.Type(logger.TypePredb).Debugf("release is too old (%s)", time.Since(pre.At).String())
			return
		}
	}

	// -- Get Category ----------------------------
	logWithRef.Debugf("getting internal category for preDBCategory: %v", pre.Category)
	cat := a.GetCategoryByName(pre.Category.Name)
	if cat == nil {
		logWithRef.Debugf("no category found for preDBCategory: %v", pre.Category.Name)
		return
	}

	// -- Check if category is allowed ---------
	if accepted, err := cat.Accepts(r.Name, r.Size); !accepted {
		logWithRef.Infof("release %s is not accepted by category %s: %s", r.Name, cat.Name, err.Error())
		return
	}

	// -- Save Release ----------------------------
	if err := r.Save(pre); err != nil {
		logWithRef.Type(logger.TypeGeneric).Errorf("failed to save release %s: %s", r.Name, err.Error())
		return
	}

	r.Source.SumReleasesDownloaded++
	if err := r.Source.Save(); err != nil {
		logWithRef.Type(logger.TypeGeneric).Errorf("failed to save source %s: %s", r.Source.Name, err.Error())
		return
	}

	a.OnReleaseAdded(r)

	// -- add to pending releases queue -----------

	a.pendingReleases.Store(r.Hash, &Release{
		UID:         r.UID,
		Name:        r.Name,
		Hash:        r.Hash,
		Size:        r.Size,
		State:       r.State,
		SourceUID:   r.Source.UID,
		Category:    string(cat.Name),
		CategoryRaw: pre.CategoryRaw,
		Pre:         pre.At,
		MetaFiles:   r.MetaFiles,
	})

}

// processPendingReleases processes all pending releases
//   - sends new releases to fileserver
//   - checks download status of downloading releases
//   - uploads release to destination when download is finished
func (a *ATUS) processPendingReleasesTask(ctx context.Context) {

	a.pendingReleases.Range(func(key, value interface{}) bool {
		r := value.(*Release)

		logWithRef := logger.Ref(logger.RefRelease, r.UID).Type(logger.TypeRelease)

		// == handle new releases =====================================================================
		if r.State == release.StateNew {

			// -- find fileserver for release -------
			fs := a.GetFileserverForRelease(FileserverAllocationMethod(config.GetString("FILESERVER__ALLOCATION_METHOD")), r)
			if fs == nil {
				logWithRef.Errorf("no fileserver available or not enough space")
				return true
			}

			logWithRef.Debugf("assigned fileserver %s (%s)", fs.Name, fs.UID)
			a.assignFileserver(r, fs)

			// -- upload meta file -----------------------
			var torrentFile []byte

			for _, mf := range r.MetaFiles {
				if mf.State != release.MetafileStateProcessed || mf.Type != release.MetafileTypeTorrent {
					continue

				}

				file, err := mf.GetFile()
				if err != nil {
					logWithRef.Errorf("failed to get torrent file: %s", err.Error())
					a.updatePendingReleaseState(r, release.StateGeneralError)
					return true
				}
				torrentFile = file
				break
			}

			if torrentFile == nil {
				logWithRef.Errorf("no torrent file found")
				a.updatePendingReleaseState(r, release.StateGeneralError)
				return true
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			resp, err := fs.Fileserver.AddTorrent(ctx, torrentFile, r.Name+".torrent", config.GetString("FILESERVER__DOWNLOAD_LABEL"))
			// retry on common errors
			// ToDo: Check for error codes / add max retries
			if err != nil {
				logWithRef.Errorf("failed to add torrent to fileserver %s (%s). Error: %s", fs.Name, fs.UID, err.Error())
				return true
			}

			// -- check returned hash ------------------
			if resp.Hash != r.Hash {
				logWithRef.Errorf("hash of torrent file does not match the returned hash from fileserver %s (%s). Expected: %s, Got: %s", fs.Name, fs.UID, r.Hash, resp.Hash)
				// mark release as broken
				a.updatePendingReleaseState(r, release.StateGeneralError)
				return true
			}

			// -- upload successful, update release state
			logWithRef.Infof("uploaded source torrent file to fileserver %s (%s)", fs.Name, fs.UID)
			a.updatePendingReleaseState(r, release.StateDownloadInit)

			fs.Fileserver.SumFilesDownloaded++

			if err := fs.Fileserver.Save(); err != nil {
				logWithRef.Type(logger.TypeGeneric).Errorf("failed to save fileserver %s (%s): %s", fs.Fileserver.Name, fs.Fileserver.UID, err.Error())
				return true
			}

			return true

		}

		// == handle downloading releases =============================================================
		if r.State == release.StateDownloadInit || r.State == release.StateDownloading {

			fs := a.GetFileserverByUID(r.FileserverUID)
			if fs == nil {
				logWithRef.Errorf("fileserver %s not found (is nil)", r.FileserverUID)
				// This happens when the fileserver goes offline during download
				// ToDo: implement retry mechanism
				// a.updatePendingReleaseState(pr, release.STATE_ERROR)
				return true
			}

			// check if fileserver is disabled
			if !fs.Fileserver.Enabled {
				logWithRef.Errorf("fileserver %s (%s) is disabled", fs.Fileserver.Name, fs.Fileserver.UID)
				return true
			}

			downloadState, err := a.GetDownloadState(fs.Fileserver.UID, r.Hash)
			if err != nil {
				// Only log error if we saw the file on the fileserver before
				if r.State == release.StateDownloading {
					logWithRef.Errorf("failed to get download state from fileserver %s (%s): %s", fs.Fileserver.Name, fs.Fileserver.UID, err.Error())
				}
				return true
			}

			// we found the file on the fileserver, update state
			if r.State == release.StateDownloadInit {
				a.updatePendingReleaseState(r, release.StateDownloading)
			}

			// -- check if download is finished
			if downloadState.Done < 100 {
				return true
			}

			// -- download finished, update release state
			a.updatePendingReleaseState(r, release.StateDownloaded)

			return true
		}

		// == handle downloaded releases ==============================================================
		if r.State == release.StateDownloaded {
			if err := a.UploadRelease(ctx, r); err != nil {
				logWithRef.Errorf(err.Error())
				return true
			}

			logWithRef.Infof("release was successfully uploaded to tracker")
		}

		return true

	})

}

func (a *ATUS) UploadRelease(ctx context.Context, r *Release) error {

	fs := a.GetFileserverByUID(r.FileserverUID)
	if fs == nil {
		return fmt.Errorf("fileserver %s not found (is nil)", r.FileserverUID)
	}

	// check if all meta files are processed
	for _, mf := range r.MetaFiles {
		if mf.State != release.MetafileStateProcessed && mf.State != release.MetafileStateError {
			return fmt.Errorf("meta file %s is not processed", mf.FileName)
		}
	}

	newDict, err := a.UploadReleaseToTracker(ctx, r)
	if err != nil {
		// There is currently no auto-retry mechanism
		// Once a release is in an error state, it will have to be manually uploaded through the web interface
		a.updatePendingReleaseState(r, release.StateUploadError)
		return fmt.Errorf("failed to upload release to tracker: %s", err.Error())
	}

	// the release was uploaded successfully
	// we now have to prepare the .torrent file with the trackers announce url
	newDict.Announce = config.GetString("UPLOAD__USER_ANNOUNCE_URL")
	newTorrent, err := newDict.BEncode()
	if err != nil {
		a.updatePendingReleaseState(r, release.StateUploadError)
		return fmt.Errorf("failed to encode torrent file: %s", err.Error())
	}

	// send new torrent to fileserver
	if _, err := fs.Fileserver.AddTorrent(ctx, newTorrent, r.Name+".torrent", config.GetString("FILESERVER__UPLOAD_LABEL")); err != nil {
		a.updatePendingReleaseState(r, release.StateUploadError)
		return fmt.Errorf("failed to add destination torrent to fileserver %s (%s): %s", fs.Fileserver.Name, fs.Fileserver.UID, err.Error())
	}

	// update release state
	a.updatePendingReleaseState(r, release.StateUploaded)

	return nil

}
