package sqlite

import (
	"database/sql"
	"log"
	"math/rand"
	"sync"
)

type uniqueIDLength int

const UIDDefaultLength uniqueIDLength = 15

// cache uids during runtime to make sure we don't return the same uid twice if not yet in db
var uidCache sync.Map

// uidInUse returns true if the given uid is in use in the given table
// column must be `uid`
func uidInUse(table string, uid string) bool {
	if _, ok := uidCache.Load(table + uid); ok {
		return true
	}

	row := Conn.QueryRow("SELECT 1 FROM `"+table+"` WHERE uid = ? LIMIT 1", uid)
	var inUse bool
	if err := row.Scan(&inUse); err != nil && err != sql.ErrNoRows {
		log.Panicf("error checking if uid %s is in use: %s", uid, err)
	}

	return inUse
}

// GenerateUID returns a unique ID not in use in the given table.
// will always return a unique value and will never return the same uid during runtime even if not entered into the db
// column must be `uid`
func GenerateUID(table string) string {
	uid := GenerateUIDLen(table, UIDDefaultLength)
	uidCache.Store(table+uid, true)
	return uid
}

// uidPool is the pool of characters to use in generating unique IDs.
// We only use url-safe characters
var uidPool = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")

// GenerateUIDLen returns a unique ID with the given length not in use in the given table.
// column must be `uid`
// @credits https://stackoverflow.com/a/22892986
func GenerateUIDLen(table string, length uniqueIDLength) string {
	poolLen := len(uidPool)
	for {

		b := make([]rune, length)
		for i := range b {
			b[i] = uidPool[rand.Intn(poolLen)]
		}

		uid := string(b)

		if !uidInUse(table, uid) {
			return uid
		}
	}
}
