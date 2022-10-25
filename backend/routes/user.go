package routes

import (
	"atus/backend/sqlite"
	"atus/backend/user"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func getSumUserAccounts() (int, error) {
	rows := sqlite.Conn.QueryRow("SELECT COUNT(1) FROM users")
	var count int

	err := rows.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func UserAuth(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(user.ContextKey).(*user.User)

	refreshToken, err := u.GenToken()
	if err != nil {
		http.Error(w, "error while generating auth token", http.StatusUnauthorized)
		return
	}

	// send refresh token
	json.NewEncoder(w).Encode(map[string]string{
		"refreshToken": refreshToken,
	})
}

func UserRegister(w http.ResponseWriter, r *http.Request) {

	var s struct {
		Username string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if s.Username == "" || s.Password == "" {
		http.Error(w, "You must provide a username and password", http.StatusBadRequest)
		return
	}

	sumUserAccounts, err := getSumUserAccounts()
	if err != nil {
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	if sumUserAccounts > 0 {
		http.Error(w, "Account setup already done. Login instead", http.StatusPreconditionFailed)
		return
	}

	u, err := user.Add(s.Username, s.Password, true)
	if err != nil {
		http.Error(w, "error while creating user", http.StatusInternalServerError)
		return
	}

	token, err := u.GenToken()
	if err != nil {
		http.Error(w, "error while generating auth token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	var s struct {
		Username string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if s.Username == "" || s.Password == "" {
		http.Error(w, "You must provide a username and password", http.StatusBadRequest)
		return
	}

	ip := GetIP(r)

	u, err := user.Login(ip, s.Username, s.Password)
	if err != nil {
		userErr := errors.New("unknown error")
		if err == sql.ErrNoRows || err == bcrypt.ErrMismatchedHashAndPassword {
			userErr = errors.New("invalid username or password")
		} else if errors.Is(err, user.ErrAuthIPBanned) {
			userErr = err
		}

		http.Error(w, userErr.Error(), http.StatusForbidden)
		return
	}

	token, err := u.GenToken()
	if err != nil {
		http.Error(w, "unknown error", http.StatusForbidden)
		panic(err)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})

}
