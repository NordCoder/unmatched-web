// Package webapp provides the static mobile UI bundled with the Go server.
package webapp

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var files embed.FS

func Static() fs.FS {
	result, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}
	return result
}
