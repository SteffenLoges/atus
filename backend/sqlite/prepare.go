package sqlite

import (
	"strings"
)

// Prepare creates tables and indices if they don't exist
func Prepare() error {

	var stmts []string

	// config
	stmts = append(stmts,
		`CREATE TABLE IF NOT EXISTS "config" (
		"name"	TEXT,
		"value"	TEXT,
		PRIMARY KEY("name")
	)`)

	stmts = append(stmts, `CREATE UNIQUE INDEX IF NOT EXISTS "config_name" ON "config" ("name")`)

	// fileservers
	stmts = append(stmts,
		`CREATE TABLE IF NOT EXISTS "fileservers" (
			"uid"	CHAR(10) NOT NULL DEFAULT NULL,
			"name"	VARCHAR(100) NOT NULL,
			"url"	TEXT NOT NULL,
			"enabled"	INTEGER NOT NULL DEFAULT 1,
			"list_interval"	INTEGER NOT NULL,
			"statistics_interval"	INTEGER NOT NULL,
			"sum_files_downloaded"	INTEGER DEFAULT 0,
			"min_free_disk_space"	INTEGER NOT NULL,
			PRIMARY KEY("uid")
		)`)

	stmts = append(stmts, `CREATE UNIQUE INDEX IF NOT EXISTS "fileservers_uid" ON "fileservers" ("uid")`)

	// log
	stmts = append(stmts,
		`CREATE TABLE IF NOT EXISTS "log" (
			"severity"	INTEGER NOT NULL,
			"type"	TEXT NOT NULL,
			"added"	TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"message"	TEXT NOT NULL,
			"release_uid"	TEXT DEFAULT ''
		)`)

	stmts = append(stmts, `CREATE INDEX IF NOT EXISTS "log_added" ON "log" ("added")`)
	stmts = append(stmts, `CREATE INDEX IF NOT EXISTS "log_brs" ON "log" ("severity", "type", "added")`)
	stmts = append(stmts, `CREATE INDEX IF NOT EXISTS "log_severity" ON "log" ("severity")`)
	stmts = append(stmts, `CREATE INDEX IF NOT EXISTS "log_type" ON "log" ("type")`)
	stmts = append(stmts, `CREATE INDEX IF NOT EXISTS "log_uidad" ON "log" ("release_uid", "added")`)

	// releases
	stmts = append(stmts,
		`CREATE TABLE IF NOT EXISTS "releases" (
			"uid"	TEXT NOT NULL UNIQUE,
			"hash"	TEXT NOT NULL,
			"name"	TEXT NOT NULL UNIQUE,
			"name_raw"	TEXT NOT NULL,
			"state"	TEXT NOT NULL,
			"pre"	TEXT NOT NULL,
			"category"	TEXT NOT NULL,
			"category_raw"	TEXT NOT NULL,
			"size"	INTEGER NOT NULL,
			"added"	TEXT NOT NULL,
			"uploaded"	TEXT,
			"source_uid"	TEXT NOT NULL,
			"fileserver_uid"	TEXT DEFAULT '',
			PRIMARY KEY("uid")
		)`)

	stmts = append(stmts, `CREATE UNIQUE INDEX IF NOT EXISTS "releases_name" ON "releases" ("name")`)

	// release_metafiles
	stmts = append(stmts,
		`CREATE TABLE IF NOT EXISTS "release_metafiles" (
			"uid"	TEXT NOT NULL UNIQUE,
			"release_uid"	TEXT NOT NULL,
			"index"	INTEGER,
			"type"	TEXT NOT NULL,
			"state"	TEXT NOT NULL,
			"file_name"	TEXT,
			"info"	TEXT DEFAULT '{}'
		)`)

	stmts = append(stmts, `CREATE INDEX IF NOT EXISTS "release_metafiles_uid" ON "release_metafiles" ("release_uid")`)

	// sources
	stmts = append(stmts,
		`CREATE TABLE IF NOT EXISTS "sources" (
			"uid"	CHAR(10) NOT NULL DEFAULT NULL,
			"name"	VARCHAR(100) NOT NULL,
			"favicon"	TEXT NOT NULL,
			"enabled"	INTEGER NOT NULL DEFAULT 0,
			"cookies"	TEXT NOT NULL,
			"rss_url"	TEXT NOT NULL,
			"rss_interval"	INTEGER NOT NULL,
			"request_waittime"	INTEGER NOT NULL,
			"last_check"	TEXT DEFAULT NULL,
			"meta_path"	TEXT NOT NULL,
			"meta_path_use_as_key"	TINYINT NOT NULL,
			"image_path"	TEXT NOT NULL,
			"image_path_use_as_key"	TINYINT NOT NULL,
			"times_checked"	INTEGER DEFAULT 0,
			"sum_torrents_downloaded"	INTEGER DEFAULT 0,
			"sum_images_downloaded"	INTEGER DEFAULT 0,
			"sum_releases_downloaded"	INTEGER DEFAULT 0,
			PRIMARY KEY("uid")
		)`)

	stmts = append(stmts, `CREATE INDEX IF NOT EXISTS "sources_name" ON "sources" ("name")`)

	// users
	stmts = append(stmts,
		`CREATE TABLE IF NOT EXISTS "users" (
			"uid"	TEXT NOT NULL,
			"is_master"	INTEGER NOT NULL DEFAULT 0,
			"name"	TEXT NOT NULL,
			"last_login"	TEXT NOT NULL DEFAULT '0000-00-00 00:00:00',
			"password_hash"	TEXT NOT NULL,
			UNIQUE("name")
		)`)

	stmts = append(stmts, `CREATE INDEX IF NOT EXISTS "users_name" ON "users" ("name")`)

	if _, err := Conn.Exec(strings.Join(stmts, ";")); err != nil {
		return err
	}

	return nil
}
