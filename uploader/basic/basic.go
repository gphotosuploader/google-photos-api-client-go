package basic

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader"
)

// BasicUploader implements an HTTP uploader.
type BasicUploader struct {
	// authenticatedClient is an HTTP client used for uploading. It needs the proper authentication in place.
	client HttpClient
	// url is the url the endpoint to upload to
	url string
	// log is a logger to send messages.
	log log.Logger
}

// HttpClient represents a HTTP client.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// UploadToken represents a pointer to the uploaded item.
type UploadToken string

// NewBasicUploader returns an Uploader or error in case of non valid configuration.
// The supplied authenticatedClient must have the proper authentication to upload files.
//
// Use WithLogger(...) and WithEndpoint(...) to customize configuration.
func NewBasicUploader(authenticatedClient HttpClient, options ...Option) (*BasicUploader, error) {
	logger := defaultLogger()
	endpoint := uploader.DefaultEndpoint

	for _, o := range options {
		switch o.Name() {
		case optkeyLogger:
			logger = o.Value().(log.Logger)
		case optkeyEndpoint:
			endpoint = o.Value().(string)
		}
	}

	return &BasicUploader{
		client: authenticatedClient,
		url:    endpoint,
		log:    logger,
	}, nil
}

// UploadFile returns the Google Photos upload token after uploading a file.
func (u BasicUploader) UploadFile(ctx context.Context, filePath string) (string, error) {
	token, err := u.upload(ctx, uploader.FileUploadItem(filePath))
	return string(token), err
}

func (u BasicUploader) upload(ctx context.Context, uploadItem uploader.UploadItem) (UploadToken, error) {
	req, err := u.prepareUploadRequest(uploadItem)
	if err != nil {
		return "", err
	}

	u.log.Debugf("Uploading %s (%d kB)", uploadItem.Name(), uploadItem.Size()/1024)
	res, err := u.client.Do(req)
	if err != nil {
		u.log.Errorf("Error while uploading %s: %s", uploadItem, err)
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		u.log.Errorf("Error while uploading %s: %s: could not read body: %s", uploadItem, res.Status, err)
		return "", err
	}
	body := string(b)

	if res.StatusCode == http.StatusOK {
		return UploadToken(body), nil
	}
	return "", fmt.Errorf("got %s: %s", res.Status, body)
}

// prepareUploadRequest returns an HTTP request to upload item.
func (u BasicUploader) prepareUploadRequest(item uploader.UploadItem) (*http.Request, error) {
	_, size, err := item.Open()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", size))
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", item.Name())
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	return req, nil
}
