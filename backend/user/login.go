package user

import (
	"atus/backend/config"
	"atus/backend/sqlite"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type failedAuthAttempt struct {
	ip         string
	username   string
	attempts   int
	bannedTill time.Time
}

var failedAuthAttempts sync.Map

var ErrAuthIPBanned = fmt.Errorf("to many failed login attempts")

func Login(ip, name, password string) (*User, error) {

	// Check if ip is banned
	if faa, ok := failedAuthAttempts.Load(ip); ok {
		f := faa.(*failedAuthAttempt)

		// reset ban if over
		if !f.bannedTill.IsZero() && f.bannedTill.Before(time.Now()) {
			failedAuthAttempts.Delete(ip)
		} else if f.bannedTill.After(time.Now()) {
			return nil, fmt.Errorf("%w. Your IP is banned until %s", ErrAuthIPBanned, f.bannedTill.Format("15:04:05"))
		}
	}

	var err error
	defer func() {
		if err == nil {
			return
		}

		// auth failed
		var faa *failedAuthAttempt
		if f, ok := failedAuthAttempts.Load(ip); ok {
			faa = f.(*failedAuthAttempt)
			faa.attempts++
		} else {
			faa = &failedAuthAttempt{
				ip:       ip,
				username: name,
				attempts: 1,
			}
		}

		if faa.attempts >= config.Base.Auth.MaxLoginAttempts {
			faa.bannedTill = time.Now().Add(config.Base.Auth.LockoutDuration)
		}

		failedAuthAttempts.Store(ip, faa)
	}()

	row := sqlite.Conn.QueryRow(
		`SELECT
			uid,
			password_hash
		FROM
			users
		WHERE
			Name = ?`, name)

	var uid string
	var passwordHash []byte
	if err = row.Scan(&uid, &passwordHash); err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password)); err != nil {
		return nil, err
	}

	// Login successful
	_, err2 := sqlite.Conn.Exec("UPDATE users SET last_login = datetime('now', 'localtime') WHERE uid = ?", uid)
	if err != nil {
		return nil, err2
	}

	return GetByUID(uid)
}
