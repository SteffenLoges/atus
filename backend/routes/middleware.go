package routes

import (
	"atus/backend/config"
	"atus/backend/helpers"
	"atus/backend/user"
	"context"
	"net/http"
	"strings"
)

func MiddlewareHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set content type to json for all responses so we don't have to do it in every handler
		// can be overwritten in the handler if needed
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Skip auth check on login/register route
		if strings.Contains(r.URL.Path, "user/login") || strings.Contains(r.URL.Path, "user/register") {
			next.ServeHTTP(w, r)
			return
		}

		// check if we have any accounts in the database
		// ToDo: This only needs to run once, so we should cache the result
		sumUserAccounts, err := getSumUserAccounts()
		if err != nil {
			http.Error(w, "unknown error", http.StatusInternalServerError)
			return
		}

		if sumUserAccounts == 0 {
			http.Error(w, "Account setup required", http.StatusPreconditionRequired)
			return
		}

		// check for auth token
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		// we allow the token to be passed in via the query string as well, this is primarily for websockets
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		if token == "" {
			http.Error(w, "no auth token provided", http.StatusUnauthorized)
			return
		}

		u, err := user.GetByToken(token)
		if err != nil {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, user.ContextKey, u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var apiAuthTypeContextKey helpers.ContextKey = "apiAuthType"

type apiAuthType string

const apiAuthTypeInternal apiAuthType = "INTERNAL"
const apiAuthTypeExternal apiAuthType = "EXTERNAL"

func MiddlewareAPIAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		isAuthorized := false
		authType := apiAuthTypeExternal

		// check for auth token
		authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if authToken == "" {
			authToken = r.URL.Query().Get("token")
		}

		if authToken == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if authToken == config.GetString("API__AUTH_TOKEN") {
			// request is sent from the tracker through the reverse proxy
			isAuthorized = true
		} else {
			// request is sent from the webinterface
			// check if the token is valid
			if _, err := user.GetByToken(authToken); err == nil {
				isAuthorized = true
				authType = apiAuthTypeInternal
			}
		}

		// if the user is not authorized, return 401
		if !isAuthorized {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), apiAuthTypeContextKey, authType)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
