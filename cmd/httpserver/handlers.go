package main

import (
	"fmt"
	"net/http"
)

var MyEmail = string([]byte{
	0x69, 0x67, 0x6f, 0x72, 0x40, 0x69, 0x67,
	0x6f, 0x72, 0x6d, 0x69, 0x63, 0x68, 0x61,
	0x6c, 0x61, 0x6b, 0x2e, 0x63, 0x6f, 0x6d,
})

var SitemapXML, RSSXML string

func (app *application) redirectToTLS(w http.ResponseWriter, r *http.Request) {
	target := fmt.Sprintf("https://%s%s", r.Host, r.URL.RequestURI())
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

func (app *application) getEmail(w http.ResponseWriter, r *http.Request) {
	resp := fmt.Sprintf(`<a href="mailto:%s">%s</a>`, MyEmail, MyEmail)

	w.Header().Set("Content-Type", "text/html")

	if _, err := fmt.Fprint(w, resp); err != nil {
		app.error(w, r, err)
	}
}

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{AllWritings: Writings, Quote: getRandomQuote()}
	app.render(w, r, http.StatusOK, "home.gohtml", data)
}

func (app *application) sitemap(w http.ResponseWriter, r *http.Request) {
	if SitemapXML == "" {
		SitemapXML = generateSitemap()
	}
	w.Header().Set("Content-Type", "application/xml")
	if _, err := fmt.Fprint(w, SitemapXML); err != nil {
		app.error(w, r, err)
	}
}

func (app *application) rss(w http.ResponseWriter, r *http.Request) {
	if RSSXML == "" {
		RSSXML = generateRSS()
	}
	w.Header().Set("Content-Type", "application/rss+xml")
	if _, err := fmt.Fprint(w, RSSXML); err != nil {
		app.error(w, r, err)
	}
}

func (app *application) writingView(w http.ResponseWriter, r *http.Request, writing *Writing) {
	data := &TemplateData{ActiveWriting: writing}
	app.render(w, r, http.StatusOK, "writing.gohtml", data)
}
