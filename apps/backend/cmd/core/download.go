package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	xdb "maestro/internal/db"
	xworker "maestro/internal/worker"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ayaviri/goutils/timer"
)

//  _   _    _    _   _ ____  _     _____ ____
// | | | |  / \  | \ | |  _ \| |   | ____|  _ \
// | |_| | / _ \ |  \| | | | | |   |  _| | |_) |
// |  _  |/ ___ \| |\  | |_| | |___| |___|  _ <
// |_| |_/_/   \_\_| \_|____/|_____|_____|_| \_\
//

func DownloadResourceHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var downloadUrls []string

	timer.WithTimer("obtaining download URLs from the job ID in request URL", func() {
		downloadUrls, err = getDownloadUrls(request.URL)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				"Failed to obtain file server download URLs from job ID: %v\n",
				err.Error(),
			),
			http.StatusInternalServerError,
		)
		return
	}

	writer.Header().Set("Content-Type", "application/zip")
	writer.Header().Set("Content-Disposition", "attachment; filename='cart.zip'")
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	timer.WithTimer("requesting/writing each file to zip", func() {
		err = requestAndWriteFilesToZip(downloadUrls, zipWriter)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Failed to request/write a file to the ZIP: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
	}
}

//  _   _ _____ _     ____  _____ ____  ____
// | | | | ____| |   |  _ \| ____|  _ \/ ___|
// | |_| |  _| | |   | |_) |  _| | |_) \___ \
// |  _  | |___| |___|  __/| |___|  _ < ___) |
// |_| |_|_____|_____|_|   |_____|_| \_\____/
//

// 1) Extracts the job ID from the given request's URL
// 2) Reads the job's response payload from the DB
// 3) Returns list of download URLs from response payload
func getDownloadUrls(originalRequestUrl *url.URL) ([]string, error) {
	var jobId string = strings.TrimPrefix(originalRequestUrl.Path, "/download/")
	var responsePayload string
	responsePayload, err = xdb.GetJobPayload(db, jobId)

	if err != nil {
		return []string{}, err
	}

	var message xworker.CheckoutCompletionResponse
	err = json.Unmarshal([]byte(responsePayload), &message)

	if err != nil {
		return []string{}, err
	}

	return message.DownloadUrls, nil
}

func requestAndWriteFilesToZip(downloadUrls []string, w *zip.Writer) error {
	for _, downloadUrl := range downloadUrls {
		var response *http.Response

		timer.WithTimer("requesting file", func() {
			response, err = http.Get(downloadUrl)
		})

		if err != nil {
			return err
		}

		timer.WithTimer("writing response body to ZIP writer", func() {
			err = writeFileToZip(downloadUrl, response, w)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// 1) Extracts the file name from the download URL
// 2) Creates a file header from the ZIP writer
// 3) Writes the file (in the response body) to the ZIP writer
func writeFileToZip(downloadUrl string, response *http.Response, w *zip.Writer) error {
	url, err := url.Parse(downloadUrl)

	if err != nil {
		return err
	}

	_, fileName := path.Split(url.Path)
	fileHeader := &zip.FileHeader{
		Name:   fileName,
		Method: zip.Deflate,
	}
	var fileWriter io.Writer
	fileWriter, err = w.CreateHeader(fileHeader)

	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, response.Body)

	return err
}
