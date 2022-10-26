package video

import (
	"fmt"
)

func (v *ProbedVideo) Extract(s *Stream, newFile string) ([]byte, error) {

	switch s.CodecType {
	case "video":
		// extract video and convert to x264
		return _exec("ffmpeg", "-y", "-hide_banner", "-loglevel", "warning", "-i", v.File, "-c:v", "libx264", "-vsync", "2", "-an", "-map", fmt.Sprintf("0:%d", s.Index), newFile)

	case "audio":
		// extract audio and convert to aac (ch2)
		return _exec("ffmpeg", "-y", "-hide_banner", "-loglevel", "warning", "-i", v.File, "-c:a", "aac", "-ac", "2", "-vn", "-map", fmt.Sprintf("0:%d", s.Index), newFile)

	case "subtitle":
		// extract subtitle
		return _exec("ffmpeg", "-y", "-hide_banner", "-loglevel", "warning", "-i", v.File, "-vn", "-an", "-map", fmt.Sprintf("0:%d", s.Index), newFile)

	default:
		return nil, fmt.Errorf("unsupported codec type: %s", s.CodecType)

	}
}

func (v *ProbedVideo) ExtractFrame(timestamp, newFile string) error {
	_, err := _exec("ffmpeg", "-y", "-hide_banner", "-loglevel", "warning", "-i", v.File, "-vsync", "2", "-ss", timestamp, "-frames:v", "1", "-update", "1", "-q:v", "2", newFile)
	return err
}
