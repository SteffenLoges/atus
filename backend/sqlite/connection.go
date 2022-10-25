package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB

func Connect(dsn string) error {
	var err error
	Conn, err = sql.Open("sqlite3", dsn)

	// @see https://github.com/mattn/go-sqlite3/issues/274
	Conn.SetMaxOpenConns(1)

	return err
}
