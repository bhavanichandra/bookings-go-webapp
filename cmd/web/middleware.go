package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf is a middleware that sets csrf token
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad is a middleware that sets session
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)

}
