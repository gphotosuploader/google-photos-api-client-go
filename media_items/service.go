package media_items

import (
	"context"
	"errors"
	"fmt"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"google.golang.org/api/googleapi"
	"net/http"
	"strconv"
)

// Config holds the configuration parameters for the service.
type Config struct {
	// Client should have all oAuth credentials in place.
	Client *http.Client
	URL    string
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

// maxItemsPerPage is the maximum number of media items to ask to the PhotosLibrary. Fewer media items might
// be returned than the specified number. See https://developers.google.com/photos/library/guides/list#pagination
const maxItemsPerPage = 100

var (
	// NullMediaItem is a zero value MediaItem.
	NullMediaItem = MediaItem{}

	ErrNotFound = errors.New("media item not found")
)

// Create creates a media item in the repository.
// By default, the media item will be added to the end of the library.
func (s Service) Create(ctx context.Context, mediaItem SimpleMediaItem) (MediaItem, error) {
	return s.CreateToAlbum(ctx, "", mediaItem)
}

// CreateMany creates one or more media items in the repository.
// By default, the media item(s) will be added to the end of the library.
func (s Service) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return s.CreateManyToAlbum(ctx, "", mediaItems)
}

// CreateToAlbum creates a media item in the repository.
// If an album id is specified, the media item is also added to the album.
// By default, the media item will be added to the end of the library or album.
func (s Service) CreateToAlbum(ctx context.Context, albumId string, mediaItem SimpleMediaItem) (MediaItem, error) {
	result, err := s.CreateManyToAlbum(ctx, albumId, []SimpleMediaItem{mediaItem})
	if err != nil {
		return NullMediaItem, err
	}
	return result[0], nil
}

// CreateManyToAlbum creates one or more media item(s) in the repository.
// If an album id is specified, the media item(s) are also added to the album.
// By default, the media item(s) will be added to the end of the library or album.
func (s Service) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
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
		return []MediaItem{}, err
	}
	mediaItemsResult := make([]MediaItem, len(result.NewMediaItemResults))
	for i, res := range result.NewMediaItemResults {
		// #54: MediaItem is populated if no errors occurred and the media item was created successfully.
		// If an error occurs, res.Status should have more data about the error.
		// Review the "[batchCreate on Google developers documentation]".
		//
		// [batchCreate on Google developers documentation]: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchCreate#NewMediaItemResult
		if res.MediaItem != nil {
			mediaItemsResult[i] = toMediaItem(res.MediaItem)
		}
	}
	return mediaItemsResult, nil
}

// Get returns the media item specified based on a given media item id.
// It will return ErrNotFound if the media item id does not exist.
func (s Service) Get(ctx context.Context, mediaItemId string) (*MediaItem, error) {
	result, err := s.photos.Get(mediaItemId).Context(ctx).Do()
	if err != nil && err.(*googleapi.Error).Code == http.StatusNotFound {
		return &NullMediaItem, ErrNotFound
	}
	if err != nil {
		return &NullMediaItem, fmt.Errorf("%s: %w", mediaItemId, err)
	}
	m := toMediaItem(result)
	return &m, nil
}

// ListByAlbum list all media items in the specified album.
func (s Service) ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error) {
	req := &photoslibrary.SearchMediaItemsRequest{
		AlbumId:  albumId,
		PageSize: maxItemsPerPage,
	}

	photosMediaItems := make([]*photoslibrary.MediaItem, 0)
	appendResultsFn := func(result *photoslibrary.SearchMediaItemsResponse) error {
		photosMediaItems = append(photosMediaItems, result.MediaItems...)
		return nil
	}

	if err := s.photos.Search(req).Pages(ctx, appendResultsFn); err != nil {
		return []MediaItem{}, err
	}

	mediaItems := make([]MediaItem, len(photosMediaItems))
	for i, item := range photosMediaItems {
		mediaItems[i] = toMediaItem(item)
	}

	return mediaItems, nil
}

// New returns a media items Google Photos client.
// The authenticatedClient should have all oAuth credentials in place.
func New(config Config) (*Service, error) {
	client := config.Client

	if client == nil {
		client = http.DefaultClient
	}

	photosClient, err := newPhotosLibraryClient(client, config.URL)
	if err != nil {
		return nil, err
	}

	service := &Service{
		photos: photosClient,
	}

	return service, nil
}

func newPhotosLibraryClient(authenticatedClient *http.Client, url string) (*photoslibrary.MediaItemsService, error) {
	s, err := photoslibrary.New(authenticatedClient)
	if err != nil {
		return nil, err
	}
	if url != "" {
		s.BasePath = url
	}
	return s.MediaItems, nil
}

// toMediaItem transforms a `photoslibrary.MediaItem` into a `MediaItem`.
func toMediaItem(item *photoslibrary.MediaItem) MediaItem {
	return MediaItem{
		ID:         item.Id,
		ProductURL: item.ProductUrl,
		BaseURL:    item.BaseUrl,
		MimeType:   item.MimeType,
		MediaMetadata: MediaMetadata{
			CreationTime: item.MediaMetadata.CreationTime,
			Width:        strconv.FormatInt(item.MediaMetadata.Width, 10),
			Height:       strconv.FormatInt(item.MediaMetadata.Height, 10),
		},
		Filename: item.Filename,
	}
}
