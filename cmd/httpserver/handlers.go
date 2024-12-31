package main

import (
	"net/http"
	"strings"
	"fmt"
	"time"
)

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.gohtml")
}

func (app *application) sitemap(w http.ResponseWriter, r *http.Request) {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	for _, writing := range Writings {
		loc := fmt.Sprintf(`<loc>%s/%s</loc>`, BaseURL, writing.Slug)
		lastmod := fmt.Sprintf(`<lastmod>%s</lastmod>`, writing.PublishedAt.Format(time.DateOnly))

		sb.WriteString(`<url>`)
		sb.WriteString(loc)
		sb.WriteString(lastmod)
		sb.WriteString(`</url>`)
	}

	sb.WriteString(`</urlset>`)

	w.Header().Set("Content-Type", "application/xml")

	if _, err := fmt.Fprint(w, sb.String()); err != nil {
		app.error(w, r, err)
	}
}

func (app *application) writingView(w http.ResponseWriter, r *http.Request, writing *Writing) {
	//slug := r.PathValue("slug")
	//
	//if slug == "" {
	//	http.NotFound(w, r)
	//	return
	//}
	//
	//entry, ok := view.PostIndex[slug]
	//if !ok {
	//	http.NotFound(w, r)
	//	return
	//}
	//
	//if err := view.Post(entry).Render(r.Context(), w); err != nil {
	//	app.error(w, r, err)
	//}
}
