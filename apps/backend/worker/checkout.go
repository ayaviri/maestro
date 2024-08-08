package main

import (
	"database/sql"
	"errors"
	"fmt"
	"maestro/internal"
	xdb "maestro/internal/db"
	xyoutube "maestro/internal/youtube"
	xytdlp "maestro/internal/ytdlp"
)

func DownloadCart(db *sql.DB, userId int64) error {
	var videos []xyoutube.Video

	internal.WithTimer("getting cart contents", func() {
		videos, err = xdb.GetItemsFromCart(db, userId)
	})

	if err != nil {
		return errors.New(
			fmt.Sprintf("Could not get items from cart: %v\n", err.Error()),
		)
	}

	if len(videos) == 0 {
		return errors.New("Cart is empty")
	}

	internal.WithTimer("downloading items from cart using yt-dlp", func() {
		err = xytdlp.DownloadVideos(videos)
	})

	if err != nil {
		return errors.New(
			fmt.Sprintf("Could not download vidoes using yt-dlp: %v\n", err.Error()),
		)
	}

	return nil
}
