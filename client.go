package gphotos

import (
	"net/http"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/photoservice"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

// Client is a client for uploading a media. photoslibrary does not provide `/v1/uploads` API so we implement here.
type Client struct {
	// Google Photos client
	service photoservice.Service
	// Uploader to upload new files to Google Photos
	uploader uploader.Uploader
	// cache to put albums cache
	cache cache.Cache
	// logger to send messages.
	log log.Logger
}

// NewClient constructs a new gphotos.Client from the provided HTTP client and the given options.
// The client is an HTTP client used for calling Google Photos. It needs the proper authentication in place.
//
// Use WithLogger(), WithCacher(), WithUploader() to customize it.
func NewClient(httpClient *http.Client, options ...Option) (*Client, error) {
	var service photoservice.Service
	var storer uploader.SessionStorer
	var upldr uploader.Uploader
	var err error
	logger := defaultLogger()
	cacher := defaultCacher()

	for _, o := range options {
		switch o.Name() {
		case optkeyPhotoService:
			service = o.Value().(photoservice.Service)
		case optkeyLogger:
			logger = o.Value().(log.Logger)
		case optkeyCacher:
			cacher = o.Value().(cache.Cache)
		case optkeySessionStorer:
			storer = o.Value().(uploader.SessionStorer)
		case optkeyUploader:
			upldr = o.Value().(uploader.Uploader)
		}
	}

	// Use GooglePhotosService by default.
	if service == nil {
		service, err = photoservice.NewGooglePhotosService(httpClient, WithLogger(logger))
		if err != nil {
			return nil, err
		}
	}

	// Use BasicUploader by default, as far as a SessionStorer has not been set.
	if upldr == nil {
		if storer == nil {
			upldr, err = uploader.NewBasicUploader(httpClient, WithLogger(logger))
			if err != nil {
				return nil, err
			}
		} else {
			upldr, err = uploader.NewResumableUploader(httpClient, storer, WithLogger(logger))
			if err != nil {
				return nil, err
			}
		}
	}

	return &Client{
		service:  service,
		uploader: upldr,
		cache:    cacher,
		log:      logger,
	}, nil
}

func defaultLogger() log.Logger {
	return &log.DiscardLogger{}
}

func defaultCacher() cache.Cache {
	return cache.NewCachitaCache()
}
