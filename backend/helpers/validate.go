package helpers

import "net/url"

func ValidateURL(s string) (*url.URL, bool) {
	url, err := url.Parse(s)
	if err != nil {
		return nil, false
	}

	return url, url.Scheme != "" && url.Host != ""
}
