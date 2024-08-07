package ytdlp

import (
	xyoutube "maestro/internal/youtube"
	"os/exec"
)

var err error

func DownloadVideos(v []xyoutube.Video) error {
	videoLinks := make([]string, 0)

	// TODO: Create a Videos utility class
	for _, video := range v {
		videoLinks = append(videoLinks, video.Link)
	}

	var args []string = append([]string{"-x"}, videoLinks...)

	var command *exec.Cmd = exec.Command("yt-dlp", args...)
	err = command.Run()

	return err
}
