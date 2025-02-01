package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/marcusgeorgievski/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w,r,err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w,r,http.StatusOK,"home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Get id and validate value
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Fetch snippet from db
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			http.NotFound(w, r) // snippet does not exist
		} else {
			app.serverError(w, r, err) // some other server error
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w,r,http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display form to create new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	location := fmt.Sprintf("/snippet/view/%d", id)

	http.Redirect(w, r, location, http.StatusSeeOther)

	// w.Header().Add("Location", location)
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Creating and saving new snippet..."))
}
