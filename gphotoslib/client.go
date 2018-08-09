package gphotoslib

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/palantir/stacktrace"

	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

const apiVersion = "v1"
const basePath = "https://photoslibrary.googleapis.com/"

// PhotosClient is a client for uploading a media.
// photoslibrary does not provide `/v1/uploads` API so we implement here.
type Client struct {
	*photoslibrary.Service
	Client *http.Client
}

// New constructs a new PhotosClient from an oauth httpClient
func NewClient(httpClient *http.Client) (photosClient *Client, err error) {
	photosLibraryClient, err := photoslibrary.New(httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{photosLibraryClient, httpClient}, nil
}

// GetUploadToken sends the media and returns the UploadToken.
func (client *Client) GetUploadToken(r io.Reader, filename string) (token string, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/uploads", basePath, apiVersion), r)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Goog-Upload-File-Name", filename)

	res, err := client.Client.Do(req)
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

// Upload actually uploads the media and activates it on google photos
func (client *Client) UploadFile(filePath string) (*photoslibrary.MediaItem, error) {
	filename := path.Base(filePath)
	log.Printf("Uploading %s", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed opening file")
	}
	defer file.Close()

	uploadToken, err := client.GetUploadToken(file, filename)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed getting uploadToken for %s", filename)
	}

	batchResponse, err := client.MediaItems.BatchCreate(&photoslibrary.BatchCreateMediaItemsRequest{
		NewMediaItems: []*photoslibrary.NewMediaItem{
			&photoslibrary.NewMediaItem{
				Description:     filename,
				SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: uploadToken},
			},
		},
	}).Do()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed adding media %s", filename)
	}

	if batchResponse == nil || len(batchResponse.NewMediaItemResults) != 1 {
		return nil, stacktrace.NewError("len(batchResults) should be 1")
	}
	result := batchResponse.NewMediaItemResults[0]
	if result.Status.Message != "OK" {
		return nil, stacktrace.NewError("status message should be OK, found: %s", result.Status.Message)
	}

	log.Printf("%s uploaded successfully as %s", filename, result.MediaItem.Id)
	return result.MediaItem, nil
}
