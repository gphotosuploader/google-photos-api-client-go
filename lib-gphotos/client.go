package gphotos

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"

	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

const apiVersion = "v1"
const basePath = "https://photoslibrary.googleapis.com/"

// PhotosClient is a client for uploading a media.
// photoslibrary does not provide `/v1/uploads` API so we implement here.
type Client struct {
	*photoslibrary.Service
	*http.Client
	token *oauth2.Token
}

// Token returns the value of the token used by the gphotos Client
// Cannot be used to set the token
func (c *Client) Token() *oauth2.Token {
	if c.token == nil {
		return nil
	}
	return &(*c.token)
}

// New constructs a new PhotosClient from an oauth httpClient
func NewClient(oauthHTTPClient *http.Client, maybeToken ...*oauth2.Token) (*Client, error) {
	var token *oauth2.Token
	switch len(maybeToken) {
	case 0:
	case 1:
		token = maybeToken[0]
	default:
		return nil, stacktrace.NewError("NewClient() parameters should have maximum 1 token")
	}

	photosService, err := photoslibrary.New(oauthHTTPClient)
	if err != nil {
		return nil, err
	}
	return &Client{photosService, oauthHTTPClient, token}, nil
}

// GetUploadToken sends the media and returns the UploadToken.
func (client *Client) GetUploadToken(r io.Reader, filename string) (token string, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/uploads", basePath, apiVersion), r)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("X-Goog-Upload-File-Name", filename)
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

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
func (client *Client) UploadFile(filePath string, pAlbumID ...string) (*photoslibrary.MediaItem, error) {
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

	uploadToken, err := client.GetUploadToken(file, filename)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed getting uploadToken for %s", filename)
	}

	retry := true
	retryCount := 0
	for retry != false {
		retry = false
		batchResponse, err := client.MediaItems.BatchCreate(&photoslibrary.BatchCreateMediaItemsRequest{
			AlbumId: albumID,
			NewMediaItems: []*photoslibrary.NewMediaItem{
				&photoslibrary.NewMediaItem{
					Description:     filename,
					SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: uploadToken},
				},
			},
		}).Do()
		if err != nil {
			// handle rate limit error by sleeping and retrying
			if strings.HasPrefix(err.Error(), "googleapi: Error 429") {
				log.Printf("Rate limit reached, sleeping for 10 seconds...")
				time.Sleep(10 * time.Second)
				retry = true
				retryCount++
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
		if result.Status.Message != "OK" {
			return nil, stacktrace.NewError("status message should be OK, found: %s", result.Status.Message)
		}

		log.Printf("%s uploaded successfully as %s", filename, result.MediaItem.Id)
		return result.MediaItem, nil
	}
	return nil, nil
}

func (client *Client) AlbumByName(name string) (album *photoslibrary.Album, found bool, err error) {
	listAlbumsResponse, err := client.Albums.List().Do()
	if err != nil {
		return nil, false, stacktrace.Propagate(err, "failed listing albums")
	}
	for _, album := range listAlbumsResponse.Albums {
		if album.Title == name {
			return album, true, nil
		}
	}
	return nil, false, nil
}

func (client *Client) GetOrCreateAlbumByName(albumName string) (*photoslibrary.Album, error) {
	// validate params
	{
		if albumName == "" {
			return nil, stacktrace.NewError("albumName can't be empty")
		}
	}

	// try to find album by name
	album, found, err := client.AlbumByName(albumName)
	if err != nil {
		return nil, err
	}
	if found && album != nil {
		return client.Albums.Get(album.Id).Do()
	}

	// else create album
	return client.Albums.Create(&photoslibrary.CreateAlbumRequest{
		Album: &photoslibrary.Album{
			Title: albumName,
		},
	}).Do()
}

// func (client *Client) UpsertAlbum(album photoslibrary.Album) (*photoslibrary.Album, error) {
// 	if album.Id != "" {
// 		getAlbumResponse, err := client.Albums.Get(album.Id).Fields()
// 		if err != nil {
// 			return nil, err
// 		}
// 		getAlbumResponse.Album
// 	}
// }
