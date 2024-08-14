package worker

type CheckoutRequestMessage struct {
	UserId string `json:"user_id"`
	JobId  string `json:"job_id"`
}

// TODO: The inclusion of the user ID this cart checkout is for
// paired with JWT instead of UUID bearer token for token integrity
// would be a nice touch to ensure no man-in-the-middle attack occurs
type CheckoutCompletionMessage struct {
	// The URLs to download each song from the cart from the worker's file server
	JobId        string   `json:"job_id"`
	DownloadUrls []string `json:"download_urls"`
}
