package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// templates holds the parsed template set, shared across handlers at runtime.
var templates *template.Template

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "base.html", data)
}

func (app *application) projectsList(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "projects.html", data)
}

func (app *application) projectPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "projects.html", data)
}
