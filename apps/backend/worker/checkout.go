package main

import (
	"database/sql"
	"errors"
	"fmt"
	xdb "maestro/internal/db"
	xworker "maestro/internal/worker"
	xyoutube "maestro/internal/youtube"
	xytdlp "maestro/internal/ytdlp"
	"path"
	"strings"

	"github.com/ayaviri/goutils/timer"
)

// 1) Reads the cart for the given user (contained in the message)
// 2) Downloads the cart items to disk
// 3) Returns a list of URLs, in which each one can download a
// an item from the cart from the file server in the worker's environment
func DownloadCart(
	db *sql.DB,
	message xworker.CheckoutRequestMessage,
) ([]string, error) {
	var videos []xyoutube.Video

	timer.WithTimer("getting items from user's cart", func() {
		videos, err = xdb.GetItemsFromCart(db, message.UserId)
	})

	if err != nil {
		return []string{}, fmt.Errorf(
			"Could not get items from cart: %v\n",
			err.Error(),
		)
	}

	if len(videos) == 0 {
		return []string{}, errors.New("Cart is empty")
	}

	var absoluteFilePaths []string

	timer.WithTimer("downloading items from cart using yt-dlp", func() {
		var cartDownloadDirectory string = path.Join(
			DOWNLOAD_DIRECTORY, message.JobId,
		)
		absoluteFilePaths, err = xytdlp.DownloadVideos(videos, cartDownloadDirectory)
	})

	if err != nil {
		return []string{}, fmt.Errorf(
			"Could not download vidoes using yt-dlp: %v\n",
			err.Error(),
		)
	}

	var downloadUrls []string = constructFileServerDownloadUrlsFromFilePaths(absoluteFilePaths)

	return downloadUrls, nil
}

// Constructs a download URL from the core server of the cart items downloaded
// through the given job (eg. jobId1234 => CORE_SERVER_ADDRESS + /download/jobId1234)
func ConstructCoreServerDownloadUrlFromJob(jobId string) string {
	return CORE_SERVER_ADDRESS + "/download/" + jobId
}

// Receives a list of absolute file paths
// and converts it to a list of URLs to the each file in the file server
// by trimming the download directory path (this is specific to the
// file server implementation) (eg. /User/foo/project/downloads/job/file.mp3 => FS_ADDRESS + /file.mp3)
func constructFileServerDownloadUrlsFromFilePaths(filePaths []string) []string {
	downloadUrls := make([]string, len(filePaths))

	for index, filePath := range filePaths {
		startingIndex := strings.Index(filePath, DOWNLOAD_DIRECTORY)
		strippedFilePath := strings.TrimPrefix(
			filePath[startingIndex:],
			DOWNLOAD_DIRECTORY+"/",
		)
		downloadUrls[index] = FS_ADDRESS + "/" + strippedFilePath
	}

	return downloadUrls
}

// 1) Creates a new job in the job table with the given ID for the given user
// 2) Updates its status to started
func CreateNewJob(db *sql.DB, jobId string, userId string) error {
	err = xdb.CreateNewJobWithId(db, userId, jobId)

	if err != nil {
		return err
	}

	err = xdb.UpdateJobStatus(db, jobId, xdb.StatusCreated)

	return err
}

// 1) Writes the given response payload for the given job in the job table
// 2) Updates its status to completed
func WriteJobCompletion(db *sql.DB, jobId string, responsePayload string) error {
	err = xdb.AddJobPayload(db, jobId, responsePayload)

	if err != nil {
		return err
	}

	err = xdb.UpdateJobStatus(db, jobId, xdb.StatusFinished)

	return err
}
