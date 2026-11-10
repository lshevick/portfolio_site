package main

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// registerRoutes wires the HTML page routes to their handlers. Each route maps
// to a template rendered by handlers.go; unknown paths hit NotFound.
func registerRoutes(r *httprouter.Router) {
	r.HandlerFunc("GET", "/", renderHome)
	r.HandlerFunc("GET", "/projects", renderProjectsList)
	r.HandlerFunc("GET", "/projects/:slug", func(w http.ResponseWriter, r *http.Request) {
		// httprouter matches :slug but forwards the ORIGINAL request path, so we
		// extract the slug ourselves and hand it to renderProject.
		slug := strings.TrimPrefix(r.URL.Path, "/projects/")
		if len(slug) > 0 && (strings.HasSuffix(slug, "/") || strings.ContainsAny(slug, "/")) {
			slug = strings.Split(slug, "/")[0]
		}
		renderProject(w, r, slug)
	})

	// NotFound is a field on the router, not a method.
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "page not found", http.StatusNotFound)
	})
}

// serveStatic intercepts /static/* requests and serves them from disk. It runs
// before the router (see buildHandler), so static assets never reach it. The
// "/static" prefix is stripped so a request for /static/css/output.css maps to
// <dir>/css/output.css rather than <dir>/static/css/output.css.
func serveStatic(dir string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/static") {
				fs := http.FileServer(http.Dir(dir))
				r2 := r.Clone(context.Background())
				r2.URL.Path = strings.TrimPrefix(r.URL.Path, "/static")
				fs.ServeHTTP(w, r2)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// recoverer turns a panic in any handler into a 500 instead of crashing the
// process. It is the outermost middleware so it wraps everything else.
func recoverer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// buildHandler assembles the final HTTP handler: an httprouter router with the
// page routes, wrapped by an alice middleware chain (recoverer -> serveStatic ->
// router). uiDir is where the HTML pages and static assets live.
func buildHandler(uiDir string) http.Handler {
	router := httprouter.New()
	registerRoutes(router)

	// The router implements http.Handler; wrap its ServeHTTP so the alice chain
	// (recoverer, serveStatic) sits in front of it.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	})

	chain := alice.New(recoverer(), serveStatic(filepath.Join(uiDir, "static")))
	return chain.Then(handler)
}
