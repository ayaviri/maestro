package ytdlp

import (
	xyoutube "maestro/internal/youtube"
	"os/exec"
	"strings"
)

var err error

// Downloads the given list of videos into the given download directory.
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
			"--audio-format",
			"mp3",
			"-P",
			downloadDirectory,
			"-o",
			outputTemplate,
			"--quiet",
			"--no-simulate",
			"--no-warnings",
			"--restrict-filenames",
			"--print",
			"after_move:filepath",
		},
		videoLinks...,
	)

	var out []byte
	out, err = exec.Command("yt-dlp", args...).Output()

	if err != nil {
		return []string{}, err
	}

	outString := strings.TrimSuffix(string(out), "\n")
	var absoluteFilePaths []string = strings.Split(outString, "\n")

	return absoluteFilePaths, err
}
