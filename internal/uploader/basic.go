package uploader

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/lestrrat-go/backoff"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

// SimpleUploader implements uploads to Google Photos service.
type BasicUploader struct {
	// client is an HTTP client used for uploading. It needs the proper authentication in place.
	client httpClient
	// url is the url the endpoint to upload to
	url string
	// log is a logger to send messages.
	log log.Logger
}

// NewBasicUploader returns an Uploader or error in case of non valid configuration.
// The supplied client must have the proper authentication to upload files.
//
// Use WithLogger(...) and WithEndpoint(...) to customize configuration.
func NewBasicUploader(client httpClient, options ...Option) (*BasicUploader, error) {
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

	// validate configuration options.
	if u.url == "" {
		return nil, errors.New("endpoint could not be empty")
	}

	return u, nil
}

// Upload returns the Google Photos upload token for an Upload object.
func (u *BasicUploader) Upload(ctx context.Context, item UploadItem) (UploadToken, error) {
	u.log.Debugf("Initiating file upload: type=non-resumable, file=%s", item.Name())

	req, err := u.prepareUploadRequest(item)
	if err != nil {
		return "", err
	}

	res, err := u.retryableDo(ctx, req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return UploadToken(string(b)), nil
}

// prepareUploadRequest returns an HTTP request to upload item.
func (u *BasicUploader) prepareUploadRequest(item UploadItem) (*http.Request, error) {
	r, _, err := item.Open()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.url, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", item.Name())
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	return req, nil
}

// retryableDo implements retries in a HTTP request call.
func (u *BasicUploader) retryableDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	b, cancel := defaultRetryPolicy.Start(ctx)
	defer cancel()
	for backoff.Continue(b) {
		res, err := u.client.Do(req)
		switch {
		case err == nil:
			return res, nil
		case IsRetryableError(err):
			u.log.Debugf("Error while uploading, retry: %s", err)
		case IsRateLimitError(err):
			u.log.Errorf("Rate limit reached.")
			return nil, errors.New("rate limit reached. wait ~30 seconds before trying again")
		default:
			return nil, err
		}
	}

	return nil, errors.New("retry over")
}
