package main

import (
	"github.com/ayaviri/goutils/fs"
	"github.com/ayaviri/goutils/timer"
)

func main() {
	timer.WithTimer("running file server", func() {
		servingDirectoryEnvvar := "STATIC_FILES_DIRECTORY"
		fs.InitialiseServer(3000, servingDirectoryEnvvar)
	})

}
