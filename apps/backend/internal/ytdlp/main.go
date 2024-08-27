package ytdlp

import (
	xyoutube "maestro/internal/youtube"
	"os/exec"
	"strings"
)

var err error

// Downloads the given list of videos into the given download directory.
// Returns a list of the written file paths (includes given download directory
// and file name)
// For example, if given "foo" as the download directory, list might contain
// [foo/bar.mp3 foo/baz.mp3]
func DownloadVideos(v []xyoutube.Video, downloadDirectory string) ([]string, error) {
	videoLinks := make([]string, 0)

	// TODO: Create a Videos utility class
	for _, video := range v {
		videoLinks = append(videoLinks, video.Link)
	}

	outputTemplate := "%(title)s_%(autonumber)s.%(ext)s"
	var args []string = append(
		[]string{
			"-x",
			"-P",
			downloadDirectory,
			"-o",
			outputTemplate,
			"--quiet",
			"--no-simulate",
			"--no-warnings",
			"--restrict-filenames",
			"--exec",
			"echo %(filename)q",
		},
		videoLinks...,
	)

	var out []byte
	out, err = exec.Command("yt-dlp", args...).Output()

	if err != nil {
		return []string{}, err
	}

	outString := strings.TrimSuffix(string(out), "\n")
	var fileNames []string = strings.Split(outString, "\n")

	return fileNames, err
}
