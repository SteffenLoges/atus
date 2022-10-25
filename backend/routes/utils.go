package routes

import (
	"net/http"
	"strings"
)

func GetIP(r *http.Request) string {

	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if r.RemoteAddr != "" {
		parsed := strings.Split(r.RemoteAddr, ":")
		if len(parsed) > 0 {
			return parsed[0]
		}
	}

	// This is a fallback, but should never happen
	return "[::1]"

}
