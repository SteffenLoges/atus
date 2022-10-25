package predb

import "regexp"

var xxxRegExps = []*regexp.Regexp{
	genericTagRegExp("xxx"),
}

type XXX struct{}

func isXXX(rlsName string) *XXX {
	for _, v := range xxxRegExps {
		if v.MatchString(rlsName) {
			return &XXX{}
		}
	}

	return nil
}
