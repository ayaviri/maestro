package main

import (
	"database/sql"
	"errors"
	"fmt"
	xdb "maestro/internal/db"
	xyoutube "maestro/internal/youtube"
	xytdlp "maestro/internal/ytdlp"

	"github.com/ayaviri/goutils/timer"
)

// Reads the cart for the given user and downloads in into the given
// directory (path from project root). Returns a list of downloaded file
// NAMES (eg. song.mp3)
func DownloadCart(
	db *sql.DB,
	userId string,
	downloadDirectory string,
) ([]string, error) {
	var videos []xyoutube.Video

	timer.WithTimer("getting cart contents", func() {
		videos, err = xdb.GetItemsFromCart(db, userId)
	})

	if err != nil {
		return []string{}, fmt.Errorf(
			"Could not get items from cart: %v\n",
			err.Error(),
		)
	}

	if len(videos) == 0 {
		return []string{}, errors.New("Cart is empty")
	}

	var fileDownloadPaths []string

	timer.WithTimer("downloading items from cart using yt-dlp", func() {
		fileDownloadPaths, err = xytdlp.DownloadVideos(videos, downloadDirectory)
	})

	if err != nil {
		return []string{}, fmt.Errorf(
			"Could not download vidoes using yt-dlp: %v\n",
			err.Error(),
		)
	}

	return fileDownloadPaths, nil
}
