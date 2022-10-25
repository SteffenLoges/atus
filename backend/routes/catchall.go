package routes

import (
	"atus/backend/config"
	"net/http"
	"os"
	"path"
)

// CatchAll serves the frontend
func CatchAll(w http.ResponseWriter, r *http.Request) {

	p := path.Join(config.Base.Folders.WWW, r.URL.Path)

	// serve index.html if file doesn't exist
	if _, err := os.Stat(p); os.IsNotExist(err) {
		p = path.Join(config.Base.Folders.WWW, "index.html")
	}

	http.ServeFile(w, r, p)
}
