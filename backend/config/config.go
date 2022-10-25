package config

import (
	"atus/backend/helpers"
	"atus/backend/logger"
	"atus/backend/sqlite"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
)

// Default values for config keys
// Is used when no value is found in the database
// supported types: bool, string, float64, int64
var configDefaults = map[string]interface{}{

	// -- Setup -----------------------------------
	"SETUP__SOURCE_ADDED":      false,
	"SETUP__FILESERVER_ADDED":  false,
	"SETUP__UPLOAD_CONFIGURED": false,

	// -- General ---------------------------------
	"API__AUTH_TOKEN": "",

	// -- Upload ----------------------------------
	"UPLOAD__USER_ID":              "0",
	"UPLOAD__USER_ANNOUNCE_URL":    "",
	"UPLOAD__TRACKER_ANNOUNCE_URL": "",
	"UPLOAD__API_URL":              "",
	"UPLOAD__CREATED_BY":           "ATUS",
	"UPLOAD__COMMENT":              "Torrent created by ATUS",

	// --------------------------------------------
	"FILESERVER__DOWNLOAD_LABEL":    "ATUS Download",
	"FILESERVER__UPLOAD_LABEL":      "ATUS Upload",
	"FILESERVER__ALLOCATION_METHOD": "MOST_FREE",

	// -- Samples ---------------------------------
	"SAMPLES__ENABLED":         true,
	"SAMPLES__SUM_SCREENSHOTS": int64(3),
	"SAMPLES__MIN_SIZE":        int64(helpers.MiB * 2),
	"SAMPLES__MAX_SIZE":        int64(helpers.MiB * 200),

	// -- Filters ---------------------------------
	"FILTERS__MAX_AGE": int64(0),

	"FILTERS__CATEGORY_MOVIE_ENABLED":  true,
	"FILTERS__CATEGORY_MOVIE_INCLUDES": "[]",
	"FILTERS__CATEGORY_MOVIE_EXCLUDES": "[]",
	"FILTERS__CATEGORY_MOVIE_MAX_SIZE": int64(0),

	"FILTERS__CATEGORY_TV_ENABLED":  true,
	"FILTERS__CATEGORY_TV_INCLUDES": "[]",
	"FILTERS__CATEGORY_TV_EXCLUDES": "[]",
	"FILTERS__CATEGORY_TV_MAX_SIZE": int64(0),

	"FILTERS__CATEGORY_DOCU_ENABLED":  true,
	"FILTERS__CATEGORY_DOCU_INCLUDES": "[]",
	"FILTERS__CATEGORY_DOCU_EXCLUDES": "[]",
	"FILTERS__CATEGORY_DOCU_MAX_SIZE": int64(0),

	"FILTERS__CATEGORY_APP_ENABLED":  true,
	"FILTERS__CATEGORY_APP_INCLUDES": "[]",
	"FILTERS__CATEGORY_APP_EXCLUDES": "[]",
	"FILTERS__CATEGORY_APP_MAX_SIZE": int64(0),

	"FILTERS__CATEGORY_GAME_ENABLED":  true,
	"FILTERS__CATEGORY_GAME_INCLUDES": "[]",
	"FILTERS__CATEGORY_GAME_EXCLUDES": "[]",
	"FILTERS__CATEGORY_GAME_MAX_SIZE": int64(0),

	"FILTERS__CATEGORY_AUDIO_ENABLED":  true,
	"FILTERS__CATEGORY_AUDIO_INCLUDES": "[]",
	"FILTERS__CATEGORY_AUDIO_EXCLUDES": "[]",
	"FILTERS__CATEGORY_AUDIO_MAX_SIZE": int64(0),

	"FILTERS__CATEGORY_EBOOK_ENABLED":  true,
	"FILTERS__CATEGORY_EBOOK_INCLUDES": "[]",
	"FILTERS__CATEGORY_EBOOK_EXCLUDES": "[]",
	"FILTERS__CATEGORY_EBOOK_MAX_SIZE": int64(0),

	"FILTERS__CATEGORY_XXX_ENABLED":  true,
	"FILTERS__CATEGORY_XXX_INCLUDES": "[]",
	"FILTERS__CATEGORY_XXX_EXCLUDES": "[]",
	"FILTERS__CATEGORY_XXX_MAX_SIZE": int64(0),

	"FILTERS__CATEGORY_UNKNOWN_ENABLED":  true,
	"FILTERS__CATEGORY_UNKNOWN_INCLUDES": "[]",
	"FILTERS__CATEGORY_UNKNOWN_EXCLUDES": "[]",
	"FILTERS__CATEGORY_UNKNOWN_MAX_SIZE": int64(0),
}

var cache = sync.Map{}

func get(name string) (string, bool) {

	if cached, ok := cache.Load(name); ok {
		return cached.(string), true
	}

	if _, ok := configDefaults[name]; !ok {
		logger.Fatalf("Trying to get config value for unknown key: %s", name)
	}

	row := sqlite.Conn.QueryRow("SELECT value FROM config WHERE name = ? LIMIT 1", name)

	var value string
	err := row.Scan(&value)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Fatalf("Error getting config value for %s: %s", name, err.Error())
		}
		return "", false
	}

	return value, true

}

func GetString(name string) string {
	if v, ok := get(name); ok {
		return v
	}

	return configDefaults[name].(string)
}

func GetInt64(name string) int64 {
	if v, ok := get(name); ok {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}

	return configDefaults[name].(int64)
}

func GetFloat64(name string) float64 {
	if v, ok := get(name); ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}

	return configDefaults[name].(float64)
}

func GetBool(name string) bool {
	if v, ok := get(name); ok {
		return v == "true"
	}

	return configDefaults[name].(bool)
}

func Set(name string, value interface{}) {

	if _, ok := configDefaults[name]; !ok {
		logger.Fatalf("Trying to set config value for unknown key: %s => %v", name, value)
	}

	value = fmt.Sprintf("%v", value)

	_, err := sqlite.Conn.Exec(
		`INSERT INTO config (name, value) 
		VALUES (?,?)
		ON CONFLICT (name) DO UPDATE SET value = ?`,
		name, value, value)

	if err != nil {
		logger.Fatalf("Error setting config value for %s: %s", name, err.Error())
	}

	cache.Store(name, value)

}
