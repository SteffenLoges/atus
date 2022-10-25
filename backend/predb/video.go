package predb

import "regexp"

type Video struct{}

var videoVarsRegExp = []*regexp.Regexp{
	// codecs
	genericTagRegExp("(tv)?xvid"),
	genericTagRegExp("divx"),
	genericTagRegExp("[x|h]26[4|5]"),

	// resolutions
	genericTagRegExp("[240|360|480|576|720|1080|2160|4320][p|i]"),

	// formats
	genericTagRegExp("s?vcd"),
	genericTagRegExp("dvd-?r(ip)?"),

	// screener
	genericTagRegExp("dvdscr(eener)?"),
	genericTagRegExp("screener"),
}

func isVideo(rlsName string) *Video {
	for _, v := range videoVarsRegExp {
		if v.MatchString(rlsName) {
			return &Video{}
		}
	}

	return nil
}
