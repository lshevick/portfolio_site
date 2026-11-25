package main

import (
	"errors"
	"net/http"

	"github.com/lshevick/internal/models"
)

// loginPage renders the sign-in form.
func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "login.html", data)
}

// loginPost verifies the submitted credentials against the database and, on
// success, records a session so later requests are treated as authenticated.
func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form struct {
		Email    string
		Password string
	}
	if err := app.decodePostForm(r, &form); err != nil {
		panic(err)
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		app.sessionManager.Put(r.Context(), "flash", models.ErrInvalidCredentials.Error())
		data := app.newTemplateData(r)
		data.Email = form.Email
		app.render(w, r, 200, "login.html", data)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "flash", "Logged in!")
	http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
}

// signupPage renders the registration form.
func (app *application) signupPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "signup.html", data)
}

// signupPost validates and inserts a new user, then logs them in.
func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form struct {
		Username string
		Email    string
		Password string
		Confirm  string
	}
	if err := app.decodePostForm(r, &form); err != nil {
		panic(err)
	}

	if form.Username == "" || form.Email == "" || form.Password == "" || form.Confirm == "" {
		app.sessionManager.Put(r.Context(), "flash", "All fields are required.")
		data := app.newTemplateData(r)
		app.render(w, r, 200, "signup.html", data)
		return
	}
	if form.Password != form.Confirm {
		app.sessionManager.Put(r.Context(), "flash", "Passwords do not match.")
		data := app.newTemplateData(r)
		app.render(w, r, 200, "signup.html", data)
		return
	}

	if err := app.users.Insert(form.Username, form.Email, form.Password); err != nil {
		var msg string
		switch {
		case errors.Is(err, models.ErrDuplicateEmails):
			msg = "That email is already registered."
		default:
			msg = err.Error()
		}
		app.sessionManager.Put(r.Context(), "flash", msg)
		data := app.newTemplateData(r)
		app.render(w, r, 200, "signup.html", data)
		return
	}

	var id int
	if err := app.users.DB.QueryRow("SELECT id FROM users WHERE email = ? LIMIT 1", form.Email).Scan(&id); err != nil {
		panic(err)
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "flash", "Account created!")
	http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
}

// dashboard is a protected landing page shown after sign-in.
func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, 200, "dashboard.html", data)
}
