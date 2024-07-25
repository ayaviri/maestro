package main

import (
    "context"
    "os"
    "fmt"
    "maestro/utils"
    "google.golang.org/api/youtube/v3"
    "google.golang.org/api/option"
)

func videosListBySearch(service *youtube.Service, parts []string, searchQuery string) {
    var call *youtube.SearchListCall = service.Search.List(parts)
    call = call.Q(searchQuery)
    var response *youtube.SearchListResponse
    var err error
    response, err = call.Do()
    utils.HandleError(err, "Unable to retrieve search list")
    var searchResults []*youtube.SearchResult = response.Items

    for _, searchResult := range searchResults {
        var snippet *youtube.SearchResultSnippet = searchResult.Snippet
        /* 
        The fields I want here are:
        - ChannelTitle
        - Description (truncate at x characters or so)
        - PublishedAt
        - Thumbnails (possibly ?)
        - Title
        
        - Link
        - Video Duration (VideosService.List -> VideoListResponse.Items[0].FileDetails.DurationMs)
        - File Size (VideosService.List -> VideoListResponse.Items[0].FileDetails.FileSize)
        - View Count (VideosService.List -> VideoListResponse.Items[0].Statistics.ViewCount)
        */
        fmt.Printf(
            "title: %v, channel title: %v\n", 
            snippet.Title, 
            snippet.ChannelTitle,
        )
    }
}

func main() {
    var service *youtube.Service
    var err error

    utils.WithTimer("creating youtube client", func() {
        ctx := context.Background()
        var apiKey string = os.Getenv("GCS_API_KEY")
        service, err = youtube.NewService(ctx, option.WithAPIKey(apiKey))
        utils.HandleError(err, "Could not create Youtube client")
    })

    var searchQuery string = "vsauce"

    utils.WithTimer(fmt.Sprintf("searching videos for %s", searchQuery), func() {
        videosListBySearch(service, []string{"snippet"}, searchQuery) 
    })
}
