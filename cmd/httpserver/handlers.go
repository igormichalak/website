package main

import (
	"net/http"
)

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.gohtml")
}

func (app *application) sitemap(w http.ResponseWriter, r *http.Request) {
	//var s strings.Builder
	//s.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	//s.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	//
	//for _, entry := range view.AllPostEntries {
	//	loc := fmt.Sprintf(`<loc>%s</loc>`, entry.URL)
	//	lastmod := fmt.Sprintf(`<lastmod>%s</lastmod>`, entry.CreatedAt.Format(time.DateOnly))
	//
	//	s.WriteString(`<url>`)
	//	s.WriteString(loc)
	//	s.WriteString(lastmod)
	//	s.WriteString(`</url>`)
	//}
	//
	//s.WriteString(`</urlset>`)
	//
	//w.Header().Set("Content-Type", "application/xml")
	//
	//if _, err := fmt.Fprint(w, s.String()); err != nil {
	//	app.error(w, r, err)
	//}
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
