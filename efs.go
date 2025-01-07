package website

import (
	"embed"
	"io/fs"
)

var StaticFS fs.FS
var TemplatesFS fs.FS
var WritingsFS fs.FS

//go:embed all:static
//go:embed view/*.gohtml view/misc/*.gohtml
//go:embed writings/*.md
var rootFS embed.FS

func init() {
	var err error

	StaticFS, err = fs.Sub(rootFS, "static")
	if err != nil {
		panic(err)
	}
	TemplatesFS, err = fs.Sub(rootFS, "view")
	if err != nil {
		panic(err)
	}
	WritingsFS, err = fs.Sub(rootFS, "writings")
	if err != nil {
		panic(err)
	}
}
