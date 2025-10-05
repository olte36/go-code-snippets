package web

import "embed"

//go:embed *.html *css
var Assets embed.FS
