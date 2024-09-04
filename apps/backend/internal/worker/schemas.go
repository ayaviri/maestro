package worker

type CheckoutRequestMessage struct {
	UserId string `json:"user_id"`
	JobId  string `json:"job_id"`
}

// The message to be sent to the client through the SSE events connection
// between the client and the core web server's /job/ endpoint
type CheckoutCompletionClientMessage struct {
	// The ID of the job that downloaded this cart. Used by the core server
	// to ensure it has picked up the correct checkout completion message
	// from the message queue
	JobId string `json:"job_id"`
	// The URL to download the cart contents (as a single zip file)
	// from the core web server
	DownloadUrl string `json:"download_url"`
}

// The payload to be written into the job table of the DB for
// the core web server to reference when downloading the cart items
// from the worker's associated file server
type CheckoutCompletionResponse struct {
	// The list of URLs corresponding to each of the items in from the
	// cart on the file server
	DownloadUrls []string `json:"download_urls"`
}
