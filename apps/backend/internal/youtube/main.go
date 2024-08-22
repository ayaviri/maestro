package youtube

import (
	"fmt"
	"html"
	"time"

	"github.com/ayaviri/goutils/timer"
	"github.com/sosodev/duration"
	"google.golang.org/api/youtube/v3"
)

var err error

//  ____   ____ _   _ _____ __  __    _    ____
// / ___| / ___| | | | ____|  \/  |  / \  / ___|
// \___ \| |   | |_| |  _| | |\/| | / _ \ \___ \
//  ___) | |___|  _  | |___| |  | |/ ___ \ ___) |
// |____/ \____|_| |_|_____|_|  |_/_/   \_\____/
//

type Video struct {
	Id              string `json:"id"`
	Title           string `json:"title"`
	ChannelTitle    string `json:"channel_title"`
	Description     string `json:"description"`
	PublishedAt     string `json:"published_at"`
	Link            string `json:"link"`
	DurationSeconds uint64 `json:"duration_seconds"`
	ViewCount       uint64 `json:"view_count"`
}

//  _   _ _____ _     ____  _____ ____  ____
// | | | | ____| |   |  _ \| ____|  _ \/ ___|
// | |_| |  _| | |   | |_) |  _| | |_) \___ \
// |  _  | |___| |___|  __/| |___|  _ < ___) |
// |_| |_|_____|_____|_|   |_____|_| \_\____/
//

func SearchVideosByQuery(
	youtubeService *youtube.Service, videoSearchQuery string,
) ([]Video, error) {
	var searchResponse *youtube.SearchListResponse

	timer.WithTimer("requesting Youtube's Search service", func() {
		parts := []string{"snippet"}
		var call *youtube.SearchListCall = youtubeService.Search.List(parts)
		call = call.Type("video")
		call = call.Q(videoSearchQuery)
		// TODO: Update once pagination parameters are introduced
		call = call.MaxResults(20)
		searchResponse, err = call.Do()
	})

	if err != nil {
		return []Video{}, err
	}

	var searchResults []*youtube.SearchResult = searchResponse.Items
	videoIds := make([]string, len(searchResults))

	timer.WithTimer("obtaining the ID for each search result", func() {
		for index, searchResult := range searchResults {
			var id *youtube.ResourceId = searchResult.Id
			videoIds[index] = id.VideoId
		}
	})

	var videosResponse *youtube.VideoListResponse

	timer.WithTimer("requesting Youtube's Video service", func() {
		parts := []string{"contentDetails", "statistics"}
		var call *youtube.VideosListCall = youtubeService.Videos.List(parts)
		call = call.Id(videoIds...)
		videosResponse, err = call.Do()
	})

	if err != nil {
		return []Video{}, err
	}

	var videoResults []*youtube.Video = videosResponse.Items
	videos := []Video{}

	timer.WithTimer(
		"mapping search and video results to desired schema, filtering out Shorts",
		func() {
			for index, searchResult := range searchResults {
				var searchSnippet *youtube.SearchResultSnippet = searchResult.Snippet
				var video *youtube.Video = videoResults[index]
				var d *duration.Duration
				d, err = duration.Parse(video.ContentDetails.Duration)

				if err != nil {
					return
				}

				var td time.Duration = d.ToTimeDuration()
				durationSeconds := uint64(td.Seconds())

				// The fact that I can't filter this out in the initial request
				// pisses me off so much
				isShort := durationSeconds <= 61

				if !isShort {
					videos = append(
						videos,
						Video{
							Id:    videoIds[index],
							Title: html.UnescapeString(searchSnippet.Title),
							ChannelTitle: html.UnescapeString(
								searchSnippet.ChannelTitle,
							),
							Description: html.UnescapeString(searchSnippet.Description),
							PublishedAt: searchSnippet.PublishedAt,
							Link: fmt.Sprintf(
								"https://www.youtube.com/watch?v=%s", videoIds[index],
							),
							DurationSeconds: durationSeconds,
							ViewCount:       video.Statistics.ViewCount,
						},
					)
				}
			}
		},
	)

	if err != nil {
		return []Video{}, err
	} else {
		return videos, nil
	}
}
