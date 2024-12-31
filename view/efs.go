package view

import "embed"

//go:embed *.gohtml misc/*.gohtml
var TemplatesFS embed.FS
