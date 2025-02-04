package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
	}

	// Inititalize buffer and execute template on it
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}


func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Decode onto target destination
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		// This is a server error
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}