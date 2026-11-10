package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// templates holds the parsed template set, shared across handlers at runtime.
var templates *template.Template

// loadTemplates parses every project page and the shared nav partial into one
// template set. ParseFiles registers each file under its BASENAME, so templates
// are looked up as "index.html", "projects.html", etc. Pages include the nav via
// {{template "nav.html" .}} (MustInclude is unavailable for these parsed files).
func loadTemplates(uiDir string) (*template.Template, error) {
	files := []string{
		filepath.Join(uiDir, "html", "partials", "nav.html"),
		filepath.Join(uiDir, "html", "pages", "index.html"),
		filepath.Join(uiDir, "html", "pages", "projects.html"),
		filepath.Join(uiDir, "html", "pages", "project-custom-keyboard.html"),
		filepath.Join(uiDir, "html", "pages", "project-project-two.html"),
		filepath.Join(uiDir, "html", "pages", "project-project-three.html"),
	}

	t := template.New("site")
	for _, f := range files {
		if _, err := t.ParseFiles(f); err != nil {
			return nil, fmt.Errorf("parse %s: %w", f, err)
		}
	}
	return t, nil
}

// projectTemplates maps a URL slug to the basename of its detail page template.
var projectTemplates = map[string]string{
	"custom-keyboard": "project-custom-keyboard.html",
	"project-two":     "project-project-two.html",
	"project-three":   "project-project-three.html",
}

// renderHome renders the landing page (ui/html/pages/index.html).
func renderHome(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"Page": "/"}
	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderProjectsList renders the projects overview grid (ui/html/pages/projects.html).
func renderProjectsList(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"Page": "/projects"}
	if err := templates.ExecuteTemplate(w, "projects.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderProject renders a single project detail page by slug.
func renderProject(w http.ResponseWriter, r *http.Request, slug string) {
	tmplName, ok := projectTemplates[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}
	data := map[string]string{"Page": "/projects/" + slug}
	if err := templates.ExecuteTemplate(w, tmplName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
