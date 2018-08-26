package gphotos

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

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

// type ClientConstructorOption func() (*Client, error)

// func FromToken(token *oauth2.Token) ClientConstructorOption {
// 	return func() (*Client, error) {
// 		httpClient := oauth2.NewClient(nil, oauth2.StaticTokenSource(token))
// 		photo NewClient(FromHTTPClient(httpClient))
// 	}
// }

// func FromHTTPClient(httpClient *http.Client,maybeToken ...*oauth2.Token) ClientConstructorOption {

// 	return func() (*Client, error) {
// 		photosService, err := photoslibrary.New(httpClient)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &Client{photosService, httpClient, nil}, nil
// 	}
// }

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
