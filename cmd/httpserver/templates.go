package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/igormichalak/website/view"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	entries, err := fs.ReadDir(view.TemplatesFS, ".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		patterns := []string{
			"misc/*.gohtml",
			name,
		}

		tmpl, err := template.New(name).ParseFS(view.TemplatesFS, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, filename string, data any) {
	tmpl, ok := app.TemplateCache[filename]
	if !ok {
		app.error(w, r, fmt.Errorf("the template %s does not exist", filename))
		return
	}

	var buf bytes.Buffer

	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		app.error(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
