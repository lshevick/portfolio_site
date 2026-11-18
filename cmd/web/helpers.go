package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	slog.Info("\n", "page", page)

	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, r, fmt.Errorf("the template %s does not exist", page))
		return
	}
	buff := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buff, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	if _, err = buff.WriteTo(w); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *application) usernameFromContext(r *http.Request) string {
	username, ok := r.Context().Value(usernameContextKey).(string)
	if !ok {
		return ""
	}
	return username
}
