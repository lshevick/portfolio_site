package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/lshevick/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	filesever := http.FileServer(http.FS(ui.Files))

	router.Handler(http.MethodGet, "/static/*filepath", filesever)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	// protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/projects", dynamic.ThenFunc(app.projectsList))
	router.Handler(http.MethodGet, "/projects/:slug", dynamic.ThenFunc(app.projectPage))

	// router.Handler(http.MethodGet, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	return standard.Then(router)
}
