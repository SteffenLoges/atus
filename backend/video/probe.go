package video

import (
	"encoding/json"
)

type ProbedVideo struct {
	File      string
	ProbeData *ProbeData
}

type ProbeData struct {
	Streams []*Stream `json:"streams"`
	Format  *Format   `json:"format"`
}

type Format struct {
	Duration string `json:"duration"`
}

type Tag struct {
	Language string `json:"language"`
	Title    string `json:"title"`
}

type Stream struct {
	Index              int    `json:"index"`
	CodecType          string `json:"codec_type"`
	CodecName          string `json:"codec_name"`
	CodecLongName      string `json:"codec_long_name"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	DisplayAspectRatio string `json:"display_aspect_ratio"`
	Tag                *Tag   `json:"tags"`
}

func Probe(file string) (*ProbedVideo, error) {

	out, err := _exec("ffprobe", "-hide_banner", "-print_format", "json", "-show_format", "-show_streams", file)
	if err != nil {
		return nil, err
	}

	var p *ProbeData
	if err := json.Unmarshal(out, &p); err != nil {
		return nil, err
	}

	return &ProbedVideo{
		File:      file,
		ProbeData: p,
	}, nil
}
