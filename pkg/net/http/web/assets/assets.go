package assets

import (
	"embed"
	"io/fs"
)

// static contains common assets used across applications contained within this
// repository.
//
//go:embed static
var assets embed.FS

// JS returns the embedded filesystem containing common JavaScript resources to
// be used across applications in this repository.
func JS() fs.FS {
	sub, err := fs.Sub(assets, "static/js")
	if err != nil {
		// this should never panic, unless the static/js directory is empty or
		// moved.
		panic(err)
	}

	return sub
}
