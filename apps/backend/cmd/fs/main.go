package main

import (
	"github.com/ayaviri/goutils/fs"
	"github.com/ayaviri/goutils/timer"
)

func main() {
	timer.WithTimer("running file server", func() {
		servingDirectoryEnvvar := "SERVING_DIRECTORY"
		fs.InitialiseServer(8001, servingDirectoryEnvvar)
	})
}
