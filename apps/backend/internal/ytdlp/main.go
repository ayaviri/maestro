package ytdlp

import (
	xyoutube "maestro/internal/youtube"
	"os/exec"
	"strings"
)

var err error

// Downloads the given list of videos nested somewhere in the given directory,
// returning the list of downloaded file NAMES (eg. song.mp3)
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
			"--print",
			outputTemplate,
			"-o",
			outputTemplate,
			"--no-simulate",
			"--no-warnings",
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
