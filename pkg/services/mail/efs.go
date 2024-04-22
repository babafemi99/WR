package mail

import "embed"

//go:embed "templates"
var EmbeddedFiles embed.FS
