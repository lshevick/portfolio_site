package main

import (
	"html/template"
	"net/http"
)

// templates holds the parsed template set, shared across handlers at runtime.
var templates *template.Template

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "home.html", data)
}

func (app *application) projectsList(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "projects.html", data)
}

func (app *application) projectCustomKeyboardPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "project-custom-keyboard.html", data)
}

func (app *application) projectTwoPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "project-project-two.html", data)
}

func (app *application) projectThreePage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "project-project-three.html", data)
}
