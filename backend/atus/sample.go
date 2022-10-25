package atus

import (
	"atus/backend/config"
	"atus/backend/logger"
	"atus/backend/release"
	"atus/backend/video"
	"fmt"
	"os"
	"path"
	"time"
)

// onNewSample is called when a new sample is found and converts it to a browser compatible format
func (a *ATUS) onNewSample(r *Release) error {

	logWithRef := logger.Ref(logger.RefRelease, r.UID)

	for _, m := range r.MetaFiles {

		// make sure the file is a sample and finished downloading
		if m.Type != release.MetafileTypeSampleVideo || m.State != release.MetafileStateDownloaded {
			continue
		}

		logWithRef.Debugf("converting sample %s", m.FileName)

		m.Info["originalFileName"] = m.FileName

		v, err := video.New(path.Join(config.Base.Folders.Data, m.ReleaseUID, m.FileName))
		if err != nil {
			logWithRef.Errorf("error probing sample video: %v", err)
			m.State = release.MetafileStateError
			return m.Save()
		}

		duration, err := time.ParseDuration(fmt.Sprintf("%ss", v.ProbeData.Format.Duration))
		if err != nil {
			logWithRef.Errorf("error parsing sample duration: %v", err)
			m.State = release.MetafileStateError
			return m.Save()
		}

		m.Info["duration"] = duration.String()

		ppv, err := v.Prepare()
		if err != nil {
			logWithRef.Errorf("error preparing sample video: %v", err)
			m.State = release.MetafileStateError
			return m.Save()
		}

		defer ppv.RemoveTempDir()

		// get file info
		for _, pv := range ppv.PreparedStreams {
			if pv.Stream.CodecType != "video" {
				continue
			}

			m.Info["width"] = fmt.Sprintf("%d", pv.Stream.Width)
			m.Info["height"] = fmt.Sprintf("%d", pv.Stream.Height)
			m.Info["codec_name"] = pv.Stream.CodecName
			m.Info["codec_long_name"] = pv.Stream.CodecLongName
			break
		}

		logWithRef.Debugf("Sample info: %s, %v", m.FileName, m.Info)

		// create screenshots
		// calc the timestamp of each screenshot dynamically
		// +2 so that the first and the last screenshot are not at the very beginning and the very end of the video
		sumScreenshots := config.GetInt64("SAMPLES__SUM_SCREENSHOTS")
		screenInterval := int64(duration/time.Second) / (sumScreenshots + 2)
		for i := int64(1); i <= sumScreenshots; i++ {

			ts := screenInterval * i
			newFileName := fmt.Sprintf("sample-screenshot-%ds-%d.jpg", ts, m.Index)
			savePath := path.Join(config.Base.Folders.Data, m.ReleaseUID, newFileName)
			err := v.ExtractFrame(fmt.Sprintf("%d", ts), savePath)
			if err != nil {
				return err
			}

			// BugFix: ffmpeg sometimes creates a 0 byte file
			// if this happens, skip this screenshot
			fileInfo, err := os.Stat(savePath)
			if err != nil {
				logWithRef.Errorf("error getting file info for %s: %v", newFileName, err)
				continue
			}

			if fileInfo.Size() == 0 {
				logWithRef.Errorf("file %s is 0 bytes, skipping", newFileName)
				continue
			}

			smv := release.NewMetaFile(m.ReleaseUID, newFileName, -1, release.MetafileTypeScreenImageFromSample, release.MetafileStateProcessed, nil, release.MetaInfo{
				"width":  m.Info["width"],
				"height": m.Info["height"],
			})

			r.MetaFiles = append(r.MetaFiles, smv)

			logWithRef.Debugf("created new sample screenshot: %s", newFileName)

			if err := smv.Save(); err != nil {
				return err
			}

			a.OnMetaFilesUpdated(r)
		}

		// -- finished processing, save
		mpdName := fmt.Sprintf("sample-manifest-%d.mpd", m.Index)
		if err := ppv.Save(path.Join(config.Base.Folders.Data, m.ReleaseUID), mpdName); err != nil {
			logWithRef.Errorf("error saving sample video: %v", err)
			m.State = release.MetafileStateError
			return m.Save()
		}

		m.Info["rawFile"] = m.FileName
		m.FileName = mpdName
		m.State = release.MetafileStateProcessed
		if err := m.Save(); err != nil {
			return err
		}

		a.OnMetaFilesUpdated(r)

		logWithRef.Infof("sample %s converted successfully", m.FileName)
	}

	return nil
}
