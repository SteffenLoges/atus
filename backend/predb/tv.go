package predb

import (
	"atus/backend/helpers"
	"regexp"
	"strconv"
	"strings"
)

type TVSeries struct {
	Name    string `json:"name"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
}

var (
	// matches:
	//   - [friends] s[1]ep[1] with extras
	//   - [24].SE[01]EP[01].DVDRiP.SVCD-FiRESToRM
	//   - [the simpsons] season[02]episode[12] cd01
	//   - [the simpsons] season[02] episode[12] cd01
	//   - [alone.the.skills.challenge].s[01]e[02].720p.web.h264-kog
	//   - [this.fool].s[01]e[08].1080p.web.h264-cakes
	//   - [fbi].s[04]e[18].german.720p.web.h264-ohd
	tvSeriesRegExp = regexp.MustCompile(`(?i)(?P<name>.+)[.\-_](s(e|eason)?)(?P<season>\d{1,3})[.\-_]?e(p|pisode)?(?P<episode>\d{1,3})[.\-_]`)

	// macthes:
	//   - [Smallville].[2]x[03].WS.HD-DSR.iNTERNAL.SVCD.SD-6
	tvSeriesRegExpAlt = regexp.MustCompile(`(?i)(?P<name>.+)[.\-_](?P<season>\d{1,3})x(?P<episode>\d{1,3})[.\-_]`)
)

// Will NOT work for releases without proper season / episode declaration
//  - CSI.307.Fight.Night.WS.HDTVRiP.SVCD-tNB
//  - Nissene.Pa.Laven.E02.NORWEiGAN.SWESUB-MDF
func isTVSeries(rlsName string, isMovie bool) *TVSeries {

	regExp := tvSeriesRegExp
	matches := regExp.FindStringSubmatch(rlsName)

	// first search returned no matches, try the alternative regexp if release is a movie
	if matches == nil && isMovie {
		regExp = tvSeriesRegExpAlt
		matches = regExp.FindStringSubmatch(rlsName)
	}

	// still no matches, return
	if matches == nil {
		return nil
	}

	// name is required
	seriesName := matches[regExp.SubexpIndex("name")]
	if seriesName == "" {
		return nil
	}

	season := 0
	seasonIndex := regExp.SubexpIndex("season")
	if seasonIndex != -1 {
		if s, err := strconv.Atoi(matches[seasonIndex]); err == nil {
			season = s
		}
	}

	episode := 0
	episodeIndex := regExp.SubexpIndex("episode")
	if episodeIndex != -1 {
		if e, err := strconv.Atoi(matches[episodeIndex]); err == nil {
			episode = e
		}
	}

	return &TVSeries{
		Name:    strings.TrimSpace(helpers.ReplaceNonAlphanumeric(seriesName, " ")),
		Season:  season,
		Episode: episode,
	}

}
