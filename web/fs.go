package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:build
var files embed.FS

func BuildHTTPFS() http.FileSystem {
	build, err := fs.Sub(files, "build")
	if err != nil {
		panic(err)
	}
	return http.FS(build)
}
