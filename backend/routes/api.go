package routes

import (
	"atus/backend/atus"
	"atus/backend/config"
	"atus/backend/release"
	"atus/backend/sqlite"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"golang.org/x/text/encoding/charmap"
)

// ToDo: allow image resizing by query params
func API__ServeFile(w http.ResponseWriter, r *http.Request) {

	requestType := r.Context().Value(apiAuthTypeContextKey).(apiAuthType)

	// check if download is requested
	isDownload := r.URL.Query().Has("download")

	// -- get file ----------------------------------------------------------------------------------

	file := path.Join(config.Base.Folders.Data, strings.TrimPrefix(r.URL.Path, "/api/data/"))

	// make sure external users can't access source torrent files
	// this should be done by the reverse proxy, but just to be sure we check it here as well
	// We allow users to access files like the source sample file. this is in case you want to
	// allow users to download the original sample file
	if requestType != apiAuthTypeInternal && strings.HasSuffix(file, ".torrent") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// check if file is servable
	stat, err := os.Stat(file)
	if err != nil || stat.IsDir() {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	buf, err := os.ReadFile(file)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// register nfo file type
	if strings.HasSuffix(file, ".nfo") {
		kind = filetype.NewType("nfo", "text/plain")
	}

	// convert nfo to utf-8 so it can be displayed in the browser
	if kind.Extension == "nfo" && !isDownload {
		buf, err = charmap.CodePage866.NewDecoder().Bytes(buf)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// set headers
	w.Header().Set("Content-Type", kind.MIME.Value)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(buf)))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Last-Modified", stat.ModTime().Format(http.TimeFormat))

	if isDownload {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path.Base(file)))
	}

	// write file
	w.Write(buf)

}

func API__Releases(a *atus.ATUS) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()

		name := query.Get("name")

		var categories []string
		if q, ok := query["category"]; ok {
			for _, c := range q {
				categories = append(categories, strings.Split(c, ",")...)
			}
		}

		var states []string
		if q, ok := query["state"]; ok {
			for _, c := range q {
				states = append(states, strings.Split(c, ",")...)
			}
		}

		orderBy := "added"
		if q := strings.ToLower(query.Get("order_by")); q != "" {
			if _, ok := map[string]bool{
				"name":     true,
				"category": true,
				"state":    true,
				"size":     true,
				"added":    true,
				"pre":      true,
			}[q]; !ok {
				http.Error(w, "invalid value for parameter 'order_by'", http.StatusBadRequest)
				return
			}
			orderBy = q
		}

		order := "desc"
		if q := strings.ToLower(query.Get("order")); q != "" {
			if _, ok := map[string]bool{
				"asc":  true,
				"desc": true,
			}[q]; !ok {
				http.Error(w, "invalid value for parameter 'order'", http.StatusBadRequest)
				return
			}
			order = q
		}

		limit := 25
		if q := query.Get("limit"); q != "" {
			n, err := strconv.Atoi(q)
			if err != nil || n < 1 || n > 100 {
				http.Error(w, "invalid value for parameter 'limit'", http.StatusBadRequest)
				return
			}
			limit = n
		}

		offset := 0
		if q := query.Get("offset"); q != "" {
			n, err := strconv.Atoi(q)
			if err != nil || n < 0 {
				http.Error(w, "invalid value for parameter 'offset'", http.StatusBadRequest)
				return
			}
			offset = n
		}

		// -- build query -------------------------------------------------------------------------------

		var clauseArr []string
		var bindings []interface{}

		if name != "" {
			clauseArr = append(clauseArr, "name LIKE ?")
			bindings = append(bindings, name)
		}

		categoriesLen := len(categories)
		if categoriesLen > 0 {
			clauseArr = append(clauseArr, fmt.Sprintf("category IN (?%s)", strings.Repeat(",?", categoriesLen-1)))
			clauseArr = append(clauseArr, "category = ?")
			for _, c := range categories {
				bindings = append(bindings, strings.ToUpper(c))
			}
		}

		statesLen := len(states)
		if statesLen > 0 {
			clauseArr = append(clauseArr, fmt.Sprintf("state IN (?%s)", strings.Repeat(",?", statesLen-1)))
			for _, s := range states {
				bindings = append(bindings, strings.ToUpper(s))
			}
		}

		var clause string
		if len(clauseArr) > 0 {
			clause += " WHERE " + strings.Join(clauseArr, " AND ")
		}

		// --------------------------------------------

		var count int
		err := sqlite.Conn.QueryRow("SELECT COUNT(1) FROM releases"+clause, bindings...).Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// --------------------------------------------

		mainQuery :=
			`SELECT 
			uid, 
			hash,
			name, 
			pre, 
			state,
			category, 
			category_raw,
			size,
			added,
			source_uid,
			fileserver_uid,
			uploaded
		FROM releases`

		mainQuery += clause + " ORDER BY " + orderBy + " " + order

		if limit > 0 {
			mainQuery += " LIMIT " + fmt.Sprintf("%d", limit)
		}

		if offset > 0 {
			mainQuery += " OFFSET " + fmt.Sprintf("%d", offset)
		}

		releaseRows, err := sqlite.Conn.Query(mainQuery, bindings...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer releaseRows.Close()

		var releases []map[string]interface{}
		for releaseRows.Next() {
			var uid, name, pre, category, categoryRaw, addedRaw, hash, sourceUID, fileserverUID string
			var uploaded sql.NullString
			var state release.ReleaseState
			var size int64

			if err := releaseRows.Scan(
				&uid,
				&hash,
				&name,
				&pre,
				&state,
				&category,
				&categoryRaw,
				&size,
				&addedRaw,
				&sourceUID,
				&fileserverUID,
				&uploaded,
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			added, err := time.Parse(time.RFC3339, addedRaw)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// sourceName := ""
			// if source := a.GetSourceByUID(sourceUID); source != nil {
			// 	sourceName = source.Source.Name
			// }

			downloadState, _ := a.GetDownloadState(fileserverUID, hash)
			// fileserverName := ""
			// if fileserver := a.GetFileserverByUID(fileserverUID); fileserver != nil {
			// 	fileserverName = fileserver.Fileserver.Name
			// }

			release := map[string]interface{}{
				"uid":         uid,
				"hash":        hash,
				"name":        name,
				"pre":         pre,
				"category":    category,
				"categoryRaw": categoryRaw,
				"size":        size,
				"added":       added,
				// "sourceName":    sourceName,
				"downloadState": downloadState,
				"state": map[string]interface{}{
					"state":      state,
					"uploadDate": uploaded.String,
				},
				// "fileserver": map[string]interface{}{
				// 	"UID":            fileserverUID,
				// 	"fileserverName": fileserverName,
				// },
			}

			releases = append(releases, release)
		}

		// for _, r := range releases {
		// 	metaFiles, _ := release.GetMetaFiles(r["uid"].(string), "")
		// 	r["metaFiles"] = metaFiles
		// }

		// I've commented out possibly sensitive data. You can uncomment it if you want to see it in the
		// response but note that if you are using the nginx config i've provided in atus-tracker-api repo,
		// all the data is accessible by anyone WITHOUT authentication

		json.NewEncoder(w).Encode(map[string]interface{}{
			"count":    count,
			"releases": releases,
		})

	}
}
