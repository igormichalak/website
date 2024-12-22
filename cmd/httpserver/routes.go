package main

import (
	"io/fs"
	"net/http"
	"os"
	"strings"
)

const StaticDirPath = "./static"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	staticFS := os.DirFS(StaticDirPath)
	fileServer := http.FileServerFS(staticFS)

	entries, err := fs.ReadDir(staticFS, ".")
	if err != nil {
		panic(err)
	}
	topLevelPaths := make([]string, len(entries))
	for i, entry := range entries {
		topLevelPaths[i] = entry.Name()
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		var firstSegment string
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) > 0 {
			firstSegment = segments[0]
		}
		switch firstSegment {
		case "":
			app.homeView(w, r)
		case "sitemap.xml":
			app.sitemap(w, r)
		default:
			for _, p := range topLevelPaths {
				if firstSegment == p {
					fileServer.ServeHTTP(w, r)
					return
				}
			}
			for i := range Writings {
				if firstSegment == Writings[i].Slug {
					app.writingView(w, r, &Writings[i])
					return
				}
			}
			http.NotFound(w, r)
		}
	})

	return app.recoverPanic(app.wwwRedirect(app.securityHeaders(app.logRequest(mux))))
}
