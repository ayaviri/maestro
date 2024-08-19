package main

import (
	"encoding/json"
	"fmt"
	xdb "maestro/internal/db"
	xhttp "maestro/internal/http"
	xyoutube "maestro/internal/youtube"
	"net/http"
	"net/url"
	"time"

	"github.com/ayaviri/goutils/timer"
)

//  ____   ____ _   _ _____ __  __    _    ____
// / ___| / ___| | | | ____|  \/  |  / \  / ___|
// \___ \| |   | |_| |  _| | |\/| | / _ \ \___ \
//  ___) | |___|  _  | |___| |  | |/ ___ \ ___) |
// |____/ \____|_| |_|_____|_|  |_/_/   \_\____/
//

type VideosResponseBody struct {
	Videos []xyoutube.Video `json:"videos"`
}

//  _   _    _    _   _ ____  _     _____ ____  ____
// | | | |  / \  | \ | |  _ \| |   | ____|  _ \/ ___|
// | |_| | / _ \ |  \| | | | | |   |  _| | |_) \___ \
// |  _  |/ ___ \| |\  | |_| | |___| |___|  _ < ___) |
// |_| |_/_/   \_\_| \_|____/|_____|_____|_| \_\____/
//

func VideosResourceHandler(writer http.ResponseWriter, request *http.Request) {
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
			fmt.Sprintf("Could not get user from bearer token: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	var queryParameters url.Values = request.URL.Query()
	var videoSearchQuery string = queryParameters.Get("q")
	// TODO: Introduce pagination parameters
	videos := make([]xyoutube.Video, 0)

	// TODO: IMPLEMENT !!!
	timer.WithTimer("checking for recent searches matching query", func() {
		videos, err = xdb.GetRecentSearchResults(db, videoSearchQuery, 15*time.Minute)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not fetch recent search results: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	if len(videos) == 0 {
		timer.WithTimer("fetching videos from Youtube Data API", func() {
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
			return
		}

		timer.WithTimer("logging results to database", func() {
			err = xdb.CreateVideos(db, videos)

			if err != nil {
				return
			}

			var searchId string
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
			return
		}
	}

	// TODO: Create a struct for this response
	var responseBody []byte
	responseBody, err = json.Marshal(VideosResponseBody{Videos: videos})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				"Could not serialise videos into JSON: %v\n", err.Error(),
			),
			http.StatusInternalServerError,
		)
		return
	}

	writer.Write(responseBody)
}
