package assets

import (
	"embed"
	_ "embed"
)

//go:embed levels
var LevelFiles embed.FS
