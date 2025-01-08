package main

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/igormichalak/website"
	"github.com/justinas/alice"
)

const StaticDirPath = "./static"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	var fileServer http.Handler
	if app.Debug {
		fileServer = http.FileServerFS(os.DirFS(StaticDirPath))
	} else {
		fileServer = http.FileServerFS(website.StaticFS)
	}

	entries, err := fs.ReadDir(website.StaticFS, ".")
	if err != nil {
		panic(err)
	}
	topLevelPaths := make([]string, len(entries))
	for i, entry := range entries {
		topLevelPaths[i] = entry.Name()
	}

	mux.HandleFunc("GET /{$}", app.homeView)
	mux.HandleFunc("GET /sitemap.xml", app.sitemap)
	mux.HandleFunc("GET /rss.xml", app.rss)
	mux.HandleFunc("GET /get-email", app.getEmail)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) >= 1 {
			for _, p := range topLevelPaths {
				if segments[0] == p {
					fileServer.ServeHTTP(w, r)
					return
				}
			}
		}
		if len(segments) == 1 {
			for i := range Writings {
				if segments[0] == Writings[i].Slug {
					app.writingView(w, r, &Writings[i])
					return
				}
			}
		}
		app.notFoundView(w, r)
	})

	chain := alice.New(app.recoverPanic)
	if !app.Debug {
		chain.Append(app.wwwRedirect)
	}
	chain.Append(app.securityHeaders)
	if !app.Debug {
		chain.Append(app.cacheHeaders)
	}
	chain.Append(app.logRequest)

	return chain.Then(mux)
}
