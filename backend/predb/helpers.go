package predb

import "regexp"

// generates a regexp for a given tag
// for example:
// 	  genericTagRegExp("1080p")
//  will match:
//    Bluey.S02E50[.1080p.]WEB.h264-SALT
//    Ninja.III.The.Domination.1984[.1080p-]BluRay.X264-KiK
//    Liar_Liar_1997[_1080p_]HDDVD_x264-hV
func genericTagRegExp(v string) *regexp.Regexp {
	return regexp.MustCompile(`(?i).+[.\-_]` + v + `[.\-_]`)
}
