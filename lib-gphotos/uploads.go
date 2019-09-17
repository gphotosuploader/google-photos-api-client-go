package gphotos

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"github.com/palantir/stacktrace"
	"google.golang.org/api/googleapi"
)

// GetUploadToken sends the media and returns the UploadToken.
func (c *Client) GetUploadToken(r io.Reader, filename string) (token string, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/uploads", basePath, apiVersion), r)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("X-Goog-Upload-File-Name", filename)
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	uploadToken := string(b)
	return uploadToken, nil
}

// UploadFile actually uploads the media and activates it on google photos
func (c *Client) UploadFile(filePath string, pAlbumID ...string) (*photoslibrary.MediaItem, error) {
	// validate parameters
	if len(pAlbumID) > 1 {
		return nil, stacktrace.NewError("parameters can't include more than one albumID'")
	}
	var albumID string
	if len(pAlbumID) == 1 {
		albumID = pAlbumID[0]
	}

	filename := path.Base(filePath)
	log.Printf("Uploading %s", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed opening file")
	}
	defer file.Close()

	uploadToken, err := c.GetUploadToken(file, filename)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed getting uploadToken for %s", filename)
	}

	retry := true
	retryCount := 0
	for retry {
		// TODO: Refactor how retries are done. We should add exponential backoff
		// https://developers.google.com/photos/library/guides/best-practices#retrying-failed-requests
		retry = false // nolint
		batchResponse, err := c.MediaItems.BatchCreate(&photoslibrary.BatchCreateMediaItemsRequest{
			AlbumId: albumID,
			NewMediaItems: []*photoslibrary.NewMediaItem{
				{
					Description:     filename,
					SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: uploadToken},
				},
			},
		}).Do()
		if err != nil {
			// handle rate limit error by sleeping and retrying
			if err.(*googleapi.Error).Code == 429 {
				after, err := strconv.ParseInt(err.(*googleapi.Error).Header.Get("Retry-After"), 10, 64)
				if err != nil || after == 0 {
					after = 10
				}
				log.Printf("Rate limit reached, sleeping for %d seconds...", after)
				time.Sleep(time.Duration(after) * time.Second)
				retry = true
				continue
			} else if retryCount < 3 {
				log.Printf("Error during upload, sleeping for 10 seconds before retrying...")
				time.Sleep(10 * time.Second)
				retry = true
				retryCount++
				continue
			}
			return nil, stacktrace.Propagate(err, "failed adding media %s", filename)
		}

		if batchResponse == nil || len(batchResponse.NewMediaItemResults) != 1 {
			return nil, stacktrace.NewError("len(batchResults) should be 1")
		}
		result := batchResponse.NewMediaItemResults[0]
		if result.Status.Message != "OK" && result.Status.Message != "Success" {
			// TODO: We should use a different field like `googleapi.ServerResponse`
			return nil, stacktrace.NewError("status message should be OK/Succecc, found: %s", result.Status.Message)
		}

		log.Printf("%s uploaded successfully as %s", filename, result.MediaItem.Id)
		return result.MediaItem, nil
	}
	return nil, nil
}

// getUploadTokenResumable sends the media and returns the UploadToken.
// The uploadURL parameter will be updated
func (c *Client) getUploadTokenResumable(r io.ReadSeeker, filename string, fileSize int64, uploadURL *string) (uploadToken string, err error) {
	var offset int64

	if uploadURL == nil {
		uploadURL = new(string)
	}

	if *uploadURL != "" {
		log.Printf("Checking status of upload URL '%s'\n", *uploadURL)
		// Query previous upload status and get offset if active
		req, err := http.NewRequest("POST", *uploadURL, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Length", "0")
		req.Header.Set("X-Goog-Upload-Command", "query")

		res, err := c.Client.Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		// Get upload status
		status := res.Header.Get("X-Goog-Upload-Status")
		log.Printf("Status of upload URL '%s' is '%s'\n", *uploadURL, status)
		if status == "active" {
			offset, err = strconv.ParseInt(res.Header.Get("X-Goog-Upload-Size-Received"), 10, 64)
			if err == nil && offset > 0 && offset < fileSize {
				// Skip already uploaded part of the file
				_, err := r.Seek(offset, io.SeekStart)
				if err != nil {
					return "", err
				}
			}
		} else {
			// Other known statuses "final" and "cancelled" are both considered an Error by the official Ruby client
			// https://github.com/googleapis/google-api-ruby-client/blob/0.30.2/lib/google/apis/core/upload.rb#L250
			// Let's restart the upload from scratch
			*uploadURL = ""
		}
	}

	if *uploadURL == "" {
		// Get new upload URL
		log.Printf("Getting new upload URL for '%s'\n", filename)
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/uploads", basePath, apiVersion), nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Length", "0")
		req.Header.Set("X-Goog-Upload-Command", "start")
		req.Header.Add("X-Goog-Upload-Content-Type", "application/octet-stream")
		req.Header.Set("X-Goog-Upload-Protocol", "resumable")
		req.Header.Set("X-Goog-Upload-Raw-Size", fmt.Sprintf("%d", fileSize))

		res, err := c.Client.Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		// Read upload url
		*uploadURL = res.Header.Get("X-Goog-Upload-URL")
	}

	contentLength := fileSize - offset
	reporter := DefaultReadProgressReporter(r, filename, fileSize, offset)
	req, err := http.NewRequest("POST", *uploadURL, &reporter)
	if err != nil {
		log.Printf("Failed to prepare request: Error '%s'\n", err)
		return "", err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", contentLength))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", offset))

	res, err := c.Client.Do(req)
	if err != nil {
		log.Printf("\nFailed to process request '%s'\n", err)
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response '%s'\n", err)
		return "", err
	}
	uploadToken = string(b)
	return uploadToken, nil
}

// UploadFileResumable actually uploads the media and activaes it on google photos
// The uploadURL parameter will be updated
func (c *Client) UploadFileResumable(filePath string, uploadURL *string, pAlbumID ...string) (*photoslibrary.MediaItem, error) {
	// validate parameters
	if len(pAlbumID) > 1 {
		return nil, stacktrace.NewError("parameters can't include more than one albumID'")
	}
	var albumID string
	if len(pAlbumID) == 1 {
		albumID = pAlbumID[0]
	}

	log.Printf("Uploading %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed opening file %s", filePath)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed getting file size %s", filePath)
	}
	fileSize := fileStat.Size()

	uploadToken, err := c.getUploadTokenResumable(file, filePath, fileSize, uploadURL)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed getting uploadToken for %s", filePath)
	}

	retry := true
	retryCount := 0
	for retry {
		// TODO: Refactor how retries are done. We should add exponential backoff
		// https://developers.google.com/photos/library/guides/best-practices#retrying-failed-requests
		retry = false // nolint
		batchResponse, err := c.MediaItems.BatchCreate(&photoslibrary.BatchCreateMediaItemsRequest{
			AlbumId: albumID,
			NewMediaItems: []*photoslibrary.NewMediaItem{
				{
					Description:     filePath,
					SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: uploadToken},
				},
			},
		}).Do()
		if err != nil {
			// handle rate limit error by sleeping and retrying
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 429 {
				after, err := strconv.ParseInt(gerr.Header.Get("Retry-After"), 10, 64)
				if err != nil || after == 0 {
					after = 60
				}
				log.Printf("Rate limit reached, sleeping for %d seconds...", after)
				time.Sleep(time.Duration(after) * time.Second)
				retry = true
				continue
			} else if retryCount < 3 {
				log.Printf("Error during upload, sleeping for 10 seconds before retrying...")
				time.Sleep(10 * time.Second)
				retry = true
				retryCount++
				continue
			}
			return nil, stacktrace.Propagate(err, "failed adding media %s", filePath)
		}

		if batchResponse == nil || len(batchResponse.NewMediaItemResults) != 1 {
			return nil, stacktrace.NewError("len(batchResults) should be 1")
		}
		result := batchResponse.NewMediaItemResults[0]
		if result.Status.Message != "OK" && result.Status.Message != "Success" {
			// TODO: We should use a different field like `googleapi.ServerResponse`
			return nil, stacktrace.NewError("status message should be OK/Success, found: %s", result.Status.Message)
		}

		// Clear uploadURL as upload was completed
		*uploadURL = ""

		log.Printf("File uploaded successfully: file=%s", filePath)
		return result.MediaItem, nil
	}
	return nil, nil
}
