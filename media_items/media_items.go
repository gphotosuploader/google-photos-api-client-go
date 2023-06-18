package media_items

import (
	"context"
	"fmt"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"google.golang.org/api/googleapi"
	"net/http"
)

// A MediaItem represents a media item (e.g. photo, video etc.) in
// Google Photos.
//
// See: https://developers.google.com/photos/library/reference/rest/v1/mediaItems.
type MediaItem struct {
	// BaseURL: [Output only] A URL to the media item's bytes.
	// This should not be used as is.
	// For example, '=w2048-h1024' will set the dimensions of a media item
	// of a type photo to have a width of 2048 px and height of 1024 px.
	BaseURL string

	// Description: [Output only] Description of the media item.
	// This is shown to the user in the item's info section in the Google
	// Photos app.
	Description string

	// Filename: [Output only] Filename of the media item.
	// This is shown to the user in the item's info section in the Google Photos app.
	Filename string

	// ID: [Output only] Identifier for the media item. This is a persistent
	// identifier that can be used to identify this media item.
	ID string

	// MediaMetadata: [Output only] Metadata related to the media item, for example,
	// the height, width or creation time.
	MediaMetadata MediaMetadata

	// MimeType: [Output only] MIME type of the media item.
	MimeType string

	// ProductURL: [Output only] Google Photos URL for the media item. This link
	// will only be available to the user if they're signed in.
	ProductURL string
}

// A MediaMetadata represents the metadata for a media item.
//
// See: https://developers.google.com/photos/library/reference/rest/v1/mediaItems.
type MediaMetadata struct {
	// CreationTime: [Output only] Time when the media item was first
	// created (not when it was uploaded to Google Photos).
	CreationTime string

	// Height: [Output only] Original height (in pixels) of the media item.
	Height int64

	// Width: [Output only] Original width (in pixels) of the media item.
	Width int64
}

// A SimpleMediaItem represents a simple media item to be created in Google
// Photos via an upload token.
//
// See: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchCreate.
type SimpleMediaItem struct {
	// UploadToken: Token identifying the media bytes which have been
	// uploaded to Google.
	UploadToken string
}

// Config holds the configuration parameters for the service.
type Config struct {
	// HTTP client used to communicate with the API.
	Client *http.Client

	// [Optional] Base URL for API requests.
	// BaseURL should always be specified with a trailing slash.
	BaseURL string

	// [Optional] User agent used when communicating with the Google Photos API.
	UserAgent string
}

// Service implements a media items Google Photos client.
type Service struct {
	photos PhotosLibraryClient
}

// PhotosLibraryClient represents a Google Photos client using `gphotosuploader/googlemirror/api/photoslibrary`.
type PhotosLibraryClient interface {
	BatchCreate(batchCreateMediaItemsRequest *photoslibrary.BatchCreateMediaItemsRequest) *photoslibrary.MediaItemsBatchCreateCall
	Get(mediaItemId string) *photoslibrary.MediaItemsGetCall
	Search(searchMediaItemsRequest *photoslibrary.SearchMediaItemsRequest) *photoslibrary.MediaItemsSearchCall
}

// Create creates one media items in a user's Google Photos library.
// By default, the media item will be added to the end of the library.
func (s *Service) Create(ctx context.Context, mediaItem SimpleMediaItem) (*MediaItem, error) {
	return s.CreateToAlbum(ctx, "", mediaItem)
}

// CreateMany creates many media items in a user's Google Photos library.
// By default, the media items will be added to the end of the library.
func (s *Service) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]*MediaItem, error) {
	return s.CreateManyToAlbum(ctx, "", mediaItems)
}

// CreateToAlbum creates one media items in a user's Google Photos library.
// If an album id is specified, the media item is also added to the album.
// By default, the media item will be added to the end of the library or album.
func (s *Service) CreateToAlbum(ctx context.Context, albumId string, mediaItem SimpleMediaItem) (*MediaItem, error) {
	result, err := s.CreateManyToAlbum(ctx, albumId, []SimpleMediaItem{mediaItem})
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

// CreateManyToAlbum creates one or more media item(s) in the repository.
// If an album id is specified, the media item(s) is also added to the album.
// By default, the media item(s) will be added to the end of the library or album.
func (s *Service) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]*MediaItem, error) {
	newMediaItems := make([]*photoslibrary.NewMediaItem, len(mediaItems))
	for i, mediaItem := range mediaItems {
		newMediaItems[i] = &photoslibrary.NewMediaItem{
			SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: mediaItem.UploadToken},
		}
	}
	req := &photoslibrary.BatchCreateMediaItemsRequest{
		AlbumId:       albumId,
		NewMediaItems: newMediaItems,
	}
	result, err := s.photos.BatchCreate(req).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	mediaItemsResult := make([]*MediaItem, len(result.NewMediaItemResults))
	for i, res := range result.NewMediaItemResults {
		// #54: MediaItem is populated if no errors occurred and the media item was
		// created successfully.
		// If an error occurs, res.Status should have more data about the error.
		// In any case, we skip failed MediaItems.
		//
		// TODO: Log a message when res.MediaItem is nil. The res.Status should have the reason.
		//
		// See: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchCreate#NewMediaItemResult.
		if res.MediaItem != nil {
			mi := toMediaItem(res.MediaItem)
			mediaItemsResult[i] = &mi
		}

	}
	return mediaItemsResult, nil
}

// Get returns the media item specified based on a given media item id.
//
// Returns [ErrMediaItemNotFound] if the media item id does not exist.
func (s *Service) Get(ctx context.Context, mediaItemId string) (*MediaItem, error) {
	result, err := s.photos.Get(mediaItemId).Context(ctx).Do()
	if err != nil && err.(*googleapi.Error).Code == http.StatusNotFound {
		return nil, fmt.Errorf("%s: %w", mediaItemId, ErrMediaItemNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mediaItemId, err)
	}
	m := toMediaItem(result)
	return &m, nil
}

// maxMediaItemsPerPage is the maximum number of media items per page.
// Fewer media items might be returned than the specified number.
//
// See https://developers.google.com/photos/library/guides/list#pagination.
const maxMediaItemsPerPage = 100

// ListByAlbum list all media items in the specified album.
func (s *Service) ListByAlbum(ctx context.Context, albumId string) ([]*MediaItem, error) {
	req := &photoslibrary.SearchMediaItemsRequest{
		AlbumId:  albumId,
		PageSize: maxMediaItemsPerPage,
	}

	photosMediaItems := make([]*photoslibrary.MediaItem, 0)
	appendResultsFn := func(result *photoslibrary.SearchMediaItemsResponse) error {
		photosMediaItems = append(photosMediaItems, result.MediaItems...)
		return nil
	}

	if err := s.photos.Search(req).Pages(ctx, appendResultsFn); err != nil {
		return nil, err
	}

	mediaItems := make([]*MediaItem, len(photosMediaItems))
	for i, item := range photosMediaItems {
		mi := toMediaItem(item)
		mediaItems[i] = &mi
	}

	return mediaItems, nil
}

// New returns a media items Google Photos service.
func New(config Config) (*Service, error) {
	s, err := photoslibrary.New(config.Client)
	if err != nil {
		return nil, err
	}

	if config.BaseURL != "" {
		s.BasePath = config.BaseURL
	}

	if config.UserAgent != "" {
		s.UserAgent = config.UserAgent
	}

	service := &Service{
		photos: s.MediaItems,
	}

	return service, nil
}

func toMediaItem(item *photoslibrary.MediaItem) MediaItem {
	return MediaItem{
		ID:         item.Id,
		ProductURL: item.ProductUrl,
		BaseURL:    item.BaseUrl,
		MimeType:   item.MimeType,
		Filename:   item.Filename,
		MediaMetadata: MediaMetadata{
			CreationTime: item.MediaMetadata.CreationTime,
			Width:        item.MediaMetadata.Width,
			Height:       item.MediaMetadata.Height,
		},
	}
}
