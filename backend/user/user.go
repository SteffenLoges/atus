package user

import (
	"atus/backend/config"
	"atus/backend/helpers"
	"atus/backend/sqlite"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ContextKey helpers.ContextKey = "user"

type User struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	IsMaster  bool   `json:"isMaster"`
	LastLogin string `json:"lastLogin"`
}

func getByRaw(clause string, args ...interface{}) (*User, error) {

	row := sqlite.Conn.QueryRow(
		`SELECT
			uid,
			name,
			is_master,
			last_login
		FROM 
			users
		WHERE
			`+clause, args...)

	u := &User{}
	if err := row.Scan(&u.UID, &u.Name, &u.IsMaster, &u.LastLogin); err != nil {
		return nil, err
	}

	return u, nil
}

func GetByUID(uid string) (*User, error) {
	return getByRaw("uid = ?", uid)
}

func GetAll() ([]*User, error) {

	rows, err := sqlite.Conn.Query(
		`SELECT
			uid,
			name,
			is_master,
			last_login
		FROM 
			users
		ORDER BY
			name ASC`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.UID, &u.Name, &u.IsMaster, &u.LastLogin); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

// master accounts can't be deleted
func Add(name, password string, isMaster bool) (*User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), config.Base.Auth.BCryptCost)
	if err != nil {
		return nil, err
	}

	uid := sqlite.GenerateUID("users")

	res, err := sqlite.Conn.Exec(
		`INSERT INTO
			users
			(uid, name, password_hash, is_master)
		VALUES
			(?,?,?,?)`, uid, name, hash, isMaster)

	if err != nil {
		return nil, err
	}

	affectedRows, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if affectedRows == 0 {
		return nil, errors.New("no entry added")
	}

	return GetByUID(uid)
}

func Delete(uid string) error {

	_, err := sqlite.Conn.Exec(
		`DELETE	FROM 
			users
		WHERE
			uid = ?
		AND
			is_master = 0`, uid)

	return err
}

func (u *User) Save() error {

	_, err := sqlite.Conn.Exec(
		`UPDATE
			users
		SET
			name = ?
		WHERE
			uid = ?`, u.Name, u.UID)

	return err
}

func (u *User) SetPassword(password string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), config.Base.Auth.BCryptCost)
	if err != nil {
		return err
	}

	res, err := sqlite.Conn.Exec(
		`UPDATE
			users
		SET
			password_hash = ?
		WHERE
			uid = ?`,
		hash,
		u.UID)

	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return errors.New("no entry updated")
	}

	return nil
}
