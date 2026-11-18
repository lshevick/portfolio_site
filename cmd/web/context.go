package main

import "net/http"

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const usernameContextKey = contextKey("username")

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
