package worker

// The schema for the JSON object received by the worker for the checkout job
type DownloadCartMessageBody struct {
	UserId int64 `json:"user_id"`
	JobId  int64 `json:"job_id"`
}
