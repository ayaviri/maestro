package worker

// The schema for the JSON object received by the worker for the checkout job
type DownloadCartMessageBody struct {
	UserId string `json:"user_id"`
	JobId  string `json:"job_id"`
}

type DownloadCartResponseBody struct {
	// The URLs to download each song from the cart from the worker's file server
	DownloadUrls []string `json:"download_urls"`
}
