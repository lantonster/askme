package ui

import (
	"embed"
)

//go:embed build
var Build embed.FS

//go:embed template
var Template embed.FS
