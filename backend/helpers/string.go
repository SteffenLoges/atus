package helpers

import "regexp"

var nonAlphanumericRegExp = regexp.MustCompile(`[^a-zA-Z0-9]`)

func ReplaceNonAlphanumeric(src, repl string) string {
	return nonAlphanumericRegExp.ReplaceAllString(src, repl)
}
