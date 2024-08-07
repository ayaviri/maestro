package main

import (
	"fmt"
	"maestro/internal"
	xdb "maestro/internal/db"
	xhttp "maestro/internal/http"
	xyoutube "maestro/internal/youtube"
	xytdlp "maestro/internal/ytdlp"
	"net/http"
)

func CheckoutResourceHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user xdb.User

	internal.WithTimer("getting user from auth bearer token", func() {
		var bearerToken string
		bearerToken, _ = xhttp.GetAuthBearerToken(request)
		user, err = xdb.GetUserFromToken(db, bearerToken)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not get user from bearer token: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	var videos []xyoutube.Video

	internal.WithTimer("getting cart contents", func() {
		videos, err = xdb.GetItemsFromCart(db, user)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not get items from cart: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	if len(videos) == 0 {
		http.Error(writer, "Cart is empty", http.StatusInternalServerError)
		return
	}

	internal.WithTimer("downloading items from cart using yt-dlp", func() {
		err = xytdlp.DownloadVideos(videos)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not download videos using yt-dlp: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}
