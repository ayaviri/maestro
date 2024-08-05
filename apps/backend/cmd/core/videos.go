package main

import (
	"encoding/json"
	"fmt"
	"maestro/internal"
	xdb "maestro/internal/db"
	xhttp "maestro/internal/http"
	xyoutube "maestro/internal/youtube"
	"net/http"
)

func videosResourceHandler(writer http.ResponseWriter, request *http.Request) {
	// TODO: Update this to direct various kinds of requests for the same resource
	// to the appropriate handler
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var bearerToken string
	bearerToken, _ = xhttp.GetAuthBearerToken(request)
	var user xdb.User
	user, err = xdb.GetUserFromToken(db, bearerToken)

	if err != nil {
		http.Error(
			writer,
			"Could not get user from bearer token",
			http.StatusInternalServerError,
		)
	}

	var queryParameters url.Values = request.URL.Query()
	var videoSearchQuery string = queryParameters.Get("q")
	// TODO: Introduce pagination parameters
	var videos []xyoutube.Video

	internal.WithTimer("fetching videos from Youtube Data API", func() {
		videos, err = xyoutube.SearchVideosByQuery(
			youtubeService, videoSearchQuery,
		)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				"Could not fetch videos from Youtube Data API: %v\n",
				err.Error(),
			),
			http.StatusInternalServerError,
		)
	}

	internal.WithTimer("logging results to database", func() {
		err = xdb.CreateVideos(db, videos)
		var searchId int64
		searchId, err = xdb.CreateSearch(db, videoSearchQuery, user.Id)

		if err != nil {
			return
		}

		err = xdb.CreateSearchResults(db, searchId, videos)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				"Could not log results to database: %v\n",
				err.Error(),
			),
			http.StatusInternalServerError,
		)
	}

	var videosJson []byte
	videosJson, err = json.Marshal(videos)

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				"Could not serialise videos into JSON: %v", err.Error(),
			),
			http.StatusInternalServerError,
		)
	}

	writer.Write(videosJson)
}
