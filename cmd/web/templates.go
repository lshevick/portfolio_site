package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/lshevick/ui"
)

type templateData struct {
	CurrentYear int
	CSRFToken   string
	Flash       string
	Username    string
	// Form            *forms.Form
	IsAuthenticated bool
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// var functions = template.FuncMap{
// 	"humanDate": humanDate,
// }

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		Username:        app.usernameFromContext(r),
		CSRFToken:       nosurf.Token(r),
	}
}
