package uploader

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

// SimpleUploader implements uploads to Google Photos service.
type BasicUploader struct {
	// HTTP Client
	client *http.Client
	// URL of the endpoint to upload to
	url string

	log log.Logger
}

// NewBasicUploader returns an Uploader using the specified client or error in case
// of non valid configuration.
// The client must have the proper permissions to upload files.
//
// Use WithLogger(...) and WithEndpoint(...) to
// customize configuration.
func NewBasicUploader(client *http.Client, options ...Option) (*BasicUploader, error) {
	logger := defaultLogger()
	endpoint := defaultEndpoint()

	for _, o := range options {
		switch o.Name() {
		case optkeyLogger:
			logger = o.Value().(log.Logger)
		case optkeyEndpoint:
			endpoint = o.Value().(string)
		}
	}

	u := &BasicUploader{
		client: client,
		url:    endpoint,
		log:    logger,
	}

	return u, nil
}

// UploadFromFile returns the Google Photos upload token for the file.
func (u *BasicUploader) UploadFromFile(ctx context.Context, filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed opening file: err=%s", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed getting file maxBytes: file=%s, err=%s", filename, err)
	}
	size := fileStat.Size()

	upload := &Upload{
		r:    file,
		name: filename,
		size: size,
	}
	return u.Upload(ctx, upload)

}

// Upload returns the Google Photos upload token for an Upload object.
func (u *BasicUploader) Upload(ctx context.Context, upload *Upload) (string, error) {
	u.log.Debugf("Initiating file upload: type=non-resumable, file=%s", upload.name)

	req, err := upload.createRawUploadRequest(u.url)
	if err != nil {
		return "", err
	}

	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	token := string(b)
	return token, nil
}

// createRawUploadRequest returns a raw (non-resumable) upload request for Google Photos.
func (u *Upload) createRawUploadRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, u.r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", path.Base(u.name))
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	return req, nil

}
