package video

import (
	"atus/backend/helpers"
	"atus/backend/logger"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func New(file string) (*ProbedVideo, error) {
	v, err := Probe(file)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type PreparedStream struct {
	FileName string
	Stream   *Stream
}

type PreparedVideo struct {
	tempDir         string
	PreparedStreams []*PreparedStream
}

func (v *ProbedVideo) Prepare() (*PreparedVideo, error) {

	// make temp folder
	tempDir, err := ioutil.TempDir("", "atus_video_conv")
	if err != nil {
		return nil, err
	}

	p := &PreparedVideo{
		tempDir: tempDir,
	}

	for _, stream := range v.ProbeData.Streams {

		var newFileName string
		if stream.CodecType == "video" || stream.CodecType == "audio" {
			newFileName = fmt.Sprintf("idx_%d.mp4", stream.Index)
		} else if stream.CodecType == "subtitle" {
			newFileName = fmt.Sprintf("idx_%d.vtt", stream.Index)
		} else {
			continue
		}

		res, err := v.Extract(stream, path.Join(tempDir, newFileName))
		logger.Debugf("Extracting stream %d(%s), file: %s; %s", stream.Index, stream.CodecType, newFileName, res)
		if err != nil {
			// ignore subtitle errors (e.g. Error initializing output stream 0:0 -- Subtitle encoding currently only possible from text to text or bitmap to bitmap)
			if stream.CodecType == "subtitle" {
				continue
			}

			return nil, err
		}

		p.PreparedStreams = append(p.PreparedStreams, &PreparedStream{
			FileName: newFileName,
			Stream:   stream,
		})

	}

	return p, nil
}

func (pv *PreparedVideo) Save(savePath, mpdName string) error {

	mp4DashArgs := []string{
		fmt.Sprintf("--mpd-name=%s", mpdName),
		"-f",
		"-o", savePath,
	}

	for _, s := range pv.PreparedStreams {

		// video and audio streams need to be fragmented
		if s.Stream.CodecType == "video" || s.Stream.CodecType == "audio" {

			outputFile := path.Join(pv.tempDir, fmt.Sprintf("f-%s", s.FileName))
			err := Fragment(path.Join(pv.tempDir, s.FileName), outputFile)
			if err != nil {
				return err
			}

			mp4DashArgs = append(mp4DashArgs, outputFile)

			continue
		}

		if s.Stream.CodecType == "subtitle" {

			// it is possible that a video file has multiple subtitle streams with the same language
			// thats why we try to use the streams title as the name
			lang := fmt.Sprintf("lang-%d", s.Stream.Index)
			if s.Stream.Tag != nil {
				if s.Stream.Tag.Title != "" {
					lang = s.Stream.Tag.Title
				} else if s.Stream.Tag.Language != "" {
					lang = s.Stream.Tag.Language
				}
			}

			mp4DashArgs = append(mp4DashArgs, fmt.Sprintf("[+format=webvtt,+language=%s]%s", helpers.ReplaceNonAlphanumeric(lang, " "), path.Join(pv.tempDir, s.FileName)))
			continue
		}
	}

	ret, err := _exec("mp4dash", mp4DashArgs...)
	logger.Debugf("mp4dash %s", ret)

	return err
}

func (pv *PreparedVideo) RemoveTempDir() error {
	return os.RemoveAll(pv.tempDir)
}
