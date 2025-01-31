package main

import (
	"html/template"
	"path/filepath"

	"github.com/marcusgeorgievski/snippetbox/internal/models"
)

type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Get slice of all page filepaths that patch pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// Create template set for each page
	for _, page := range pages {
		// Extracts file name (last segment) from full filepath
		name := filepath.Base(page)

		// Parse base template file
		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Parse and add partials to existing set with base
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Parse current page template to template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}