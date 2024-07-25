package main

import (
    "time"
	"context"
	"fmt"
    "encoding/json"
	"maestro/utils"
	xhttp "maestro/utils/http"
	"net/http"
	"net/url"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
    "github.com/sosodev/duration"
)

// __   _______    ____ _     ___ _____ _   _ _____      _   _   _ _____ _   _
// \ \ / /_   _|  / ___| |   |_ _| ____| \ | |_   _|    / \ | | | |_   _| | | |
//  \ V /  | |   | |   | |    | ||  _| |  \| | | |     / _ \| | | | | | | |_| |
//   | |   | |   | |___| |___ | || |___| |\  | | |    / ___ \ |_| | | | |  _  |
//   |_|   |_|    \____|_____|___|_____|_| \_| |_|   /_/   \_\___/  |_| |_| |_|
//

var youtubeService *youtube.Service
var err error

func init() {
    utils.WithTimer("creating youtube client", func() {
        ctx := context.Background()
        var apiKey string = os.Getenv("GCS_API_KEY")
        youtubeService, err = youtube.NewService(ctx, option.WithAPIKey(apiKey))
        utils.HandleError(err, "Could not create Youtube client")
    })
}

// _____ _   _ ____  ____   ___ ___ _   _ _____ ____  
//| ____| \ | |  _ \|  _ \ / _ \_ _| \ | |_   _/ ___| 
//|  _| |  \| | | | | |_) | | | | ||  \| | | | \___ \ 
//| |___| |\  | |_| |  __/| |_| | || |\  | | |  ___) |
//|_____|_| \_|____/|_|    \___/___|_| \_| |_| |____/ 
//                                                    

func getVideosHandler(writer http.ResponseWriter, request *http.Request) {
    if (request.Method != http.MethodGet) {
        http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
        return 
    }

    xhttp.LogRequestReceipt(request)
    var queryParameters url.Values = request.URL.Query()
    var videoSearchQuery string = queryParameters.Get("q")
    // TODO: Introduce pagination parameters
    var videos []Video

    utils.WithTimer("fetching videos from Youtube Data API", func() {
        videos, err = SearchVideosByQuery(
            youtubeService, videoSearchQuery,
        )
    })

    if err != nil {
        http.Error(
            writer, 
            fmt.Sprintf(
                "Could not fetch videos from Youtube Data API: %v", err.Error(),
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
    xhttp.LogRequestCompletion(request)
}

func getHealthHandler(writer http.ResponseWriter, request *http.Request) {
    if (request.Method != http.MethodGet) {
        http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
        return 
    }

    writer.Write([]byte("maestro"))
    xhttp.LogRequestCompletion(request)
}

func initialiseServer() {
    http.HandleFunc("/health", getHealthHandler)
    http.HandleFunc("/videos", getVideosHandler)
    http.ListenAndServe(":8000", nil)
}


//  ____   ____ _   _ _____ __  __    _    ____  
// / ___| / ___| | | | ____|  \/  |  / \  / ___| 
// \___ \| |   | |_| |  _| | |\/| | / _ \ \___ \ 
//  ___) | |___|  _  | |___| |  | |/ ___ \ ___) |
// |____/ \____|_| |_|_____|_|  |_/_/   \_\____/ 
//                                               


type Video struct {
    Title string `json:"title"`
    ChannelTitle string `json:"channel_title"`
    Description string `json:"description"`
    PublishedAt  string `json:"published_at"`
    Link string `json:"link"`
    DurationSeconds uint64 `json:"duration_seconds"`
    ViewCount uint64 `json:"view_count"`
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

    utils.WithTimer("requesting Youtube's Search service", func() {
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

    utils.WithTimer("obtaining the ID for each search result", func() {
        for index, searchResult := range searchResults {
            var id *youtube.ResourceId = searchResult.Id
            videoIds[index] = id.VideoId
        }
    })

    var videosResponse *youtube.VideoListResponse

    utils.WithTimer("requesting Youtube's Video service", func() {
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

    utils.WithTimer(
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

                if (!isShort) {
                    videos = append(
                        videos, 
                        Video{
                            Title: searchSnippet.Title, 
                            ChannelTitle: searchSnippet.ChannelTitle,
                            Description: searchSnippet.Description,
                            PublishedAt: searchSnippet.PublishedAt,
                            Link: fmt.Sprintf(
                                "https://www.youtube.com/watch?v=%s", videoIds[index],
                            ),
                            DurationSeconds: durationSeconds,
                            ViewCount: video.Statistics.ViewCount,
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


//  __  __    _    ___ _   _ 
// |  \/  |  / \  |_ _| \ | |
// | |\/| | / _ \  | ||  \| |
// | |  | |/ ___ \ | || |\  |
// |_|  |_/_/   \_\___|_| \_|
//                           

func main() {
    utils.WithTimer("running HTTP server", initialiseServer)
}
