package gphotos

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"golang.org/x/xerrors"
	"google.golang.org/api/googleapi"
)

type FileItem struct {
	reader io.ReadSeeker
	size   int64
	name   string
}

// UploadFileResumable actually uploads the media and activate it on google photos
// The uploadURL parameter will be updated
func (c *Client) UploadFileResumable(filePath string, uploadURL *string, pAlbumID ...string) (*photoslibrary.MediaItem, error) {
	ctx := context.TODO() // TODO: ctx should be received (breaking change)

	// validate parameters
	if len(pAlbumID) > 1 {
		return nil, xerrors.New("parameters can't include more than one albumID")
	}
	var albumID string
	if len(pAlbumID) == 1 {
		albumID = pAlbumID[0]
	}

	c.log.Printf("[DEBUG] Initiating upload and media item creation: file=%s, type=resumable", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, xerrors.Errorf("failed opening file:  file=%s, err=%w", filePath, err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed getting file size: file=%s, err=%w", filePath, err)
	}
	fileSize := fileStat.Size()

	uploadToken, err := c.getUploadTokenResumable(file, filePath, fileSize, uploadURL)
	if err != nil {
		return nil, xerrors.Errorf("failed getting uploadToken: file=%s, err=%w", filePath, err)
	}

	c.log.Printf("[DEBUG] File has been uploaded: file=%s", filePath)

	mediaItem, err := c.createMediaItemFromUploadToken(ctx, uploadToken, albumID, filePath)
	if err != nil {
		c.log.Printf("[ERR] Failed to create media item: file=%s, err=%s", filePath, err)
		return nil, xerrors.Errorf("Error while trying to create this media item, err=%s", err)
	}

	// Clear uploadURL as upload was completed
	*uploadURL = ""

	c.log.Printf("File uploaded and media item created successfully: file=%s", filePath)
	return mediaItem, nil
}

// getUploadTokenResumable sends the media and returns the UploadToken.
// The uploadURL parameter will be updated
func (c *Client) getUploadTokenResumable(r io.ReadSeeker, filename string, fileSize int64, uploadURL *string) (uploadToken string, err error) {
	var offset int64

	if uploadURL == nil {
		uploadURL = new(string)
	}

	c.log.Printf("[DEBUG] Initiating file upload: type=`resumable`, file=%s", filename)

	if *uploadURL != "" {
		c.log.Printf("[DEBUG] Checking previous upload session: file=%s, url=%s", filename, *uploadURL)
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
		c.log.Printf("[DEBUG] Found previous upload session: file=%s, status=%s", filename, status)
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
		c.log.Printf("[DEBUG] Initiating upload session: file=%s", filename)
		url := googleapi.ResolveRelative(c.Service.BasePath, "v1/uploads")
		req, err := http.NewRequest("POST", url, nil)
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

	c.log.Printf("[DEBUG] Uploading data to Google Photos: file=%s, url=%s", filename, *uploadURL)

	contentLength := fileSize - offset
	reporter := DefaultReadProgressReporter(r, filename, fileSize, offset)
	req, err := http.NewRequest("POST", *uploadURL, &reporter)
	if err != nil {
		c.log.Printf("[ERR] Failed to prepare request: err=%s", err)
		return "", err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", contentLength))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", offset))

	res, err := c.Client.Do(req)
	if err != nil {
		c.log.Printf("[ERR] Failed to process request: err=%s", err)
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.log.Printf("[ERR] Failed to read response %s", err)
		return "", err
	}
	uploadToken = string(b)

	c.log.Printf("[DEBUG] Returning upload token: file=%s, token=%s", filename, uploadToken)

	return uploadToken, nil
}

// GetUploadToken sends the media and returns the UploadToken.
func (c *Client) GetUploadToken(r io.Reader, filename string) (token string, err error) {
	url := googleapi.ResolveRelative(c.Service.BasePath, "v1/uploads")
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return "", err
	}

	filename = path.Base(filename)
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
	ctx := context.TODO() // TODO: ctx should be received (breaking change)

	// validate parameters
	if len(pAlbumID) > 1 {
		return nil, xerrors.New("parameters can't include more than one albumID'")
	}
	var albumID string
	if len(pAlbumID) == 1 {
		albumID = pAlbumID[0]
	}

	c.log.Printf("[DEBUG] Initiating upload and media item creation: file=%s, type=not-resumable", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, xerrors.Errorf("failed opening file: err=%w", err)
	}
	defer file.Close()

	uploadToken, err := c.GetUploadToken(file, filePath)
	if err != nil {
		return nil, xerrors.Errorf("failed getting uploadToken for %s: err=%w", filePath, err)
	}

	c.log.Printf("[DEBUG] File has been uploaded: file=%s", filePath)

	mediaItem, err := c.createMediaItemFromUploadToken(ctx, uploadToken, albumID, filePath)
	if err != nil {
		c.log.Printf("[ERR] Failed to create media item: file=%s, err=%s", filePath, err)
		return nil, xerrors.Errorf("Error while trying to create this media item, err=%s", err)
	}

	c.log.Printf("File uploaded and media item created successfully: file=%s", filePath)
	return mediaItem, nil
}

func (c *Client) createMediaItemFromUploadToken(ctx context.Context, uploadToken, albumID, filename string) (*photoslibrary.MediaItem, error) {
	req := photoslibrary.BatchCreateMediaItemsRequest{
		AlbumId: albumID,
		NewMediaItems: []*photoslibrary.NewMediaItem{
			{
				Description:     filename,
				SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: uploadToken},
			},
		},
	}

	res, err := c.retryableMediaItemBatchCreateDo(ctx, &req, filename)
	if err != nil {
		return nil, err
	}

	if res == nil || len(res.NewMediaItemResults) != 1 {
		return nil, xerrors.New("len(batchResults) should be 1")
	}

	result := res.NewMediaItemResults[0]
	// `result.Status.Code` has the GRPC code returned by Google Photos API. Values can be obtained at
	// https://godoc.org/google.golang.org/genproto/googleapis/rpc/code
	if result.Status.Code != 0 {
		return nil, xerrors.New(result.Status.Message)
	}
	return result.MediaItem, nil
}
