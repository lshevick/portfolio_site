package ui

import (
	"embed"
)

//go:embed "html/base.html" "all:html" "all:static"
var Files embed.FS
