package media_items

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"google.golang.org/api/googleapi"
)

// PhotosLibraryClient represents a media items service using `gphotosuploader/googlemirror/api/photoslibrary`.
type PhotosLibraryClient interface {
	BatchCreate(batchCreateMediaItemsRequest *photoslibrary.BatchCreateMediaItemsRequest) *photoslibrary.MediaItemsBatchCreateCall
	Get(mediaItemId string) *photoslibrary.MediaItemsGetCall
	Search(searchMediaItemsRequest *photoslibrary.SearchMediaItemsRequest) *photoslibrary.MediaItemsSearchCall
}

// PhotosLibraryMediaItemsRepository represents a media items Google Photos repository.
type PhotosLibraryMediaItemsRepository struct {
	service  PhotosLibraryClient
	basePath string
}

// NewPhotosLibraryClient returns a Repository using PhotosLibrary service.
func NewPhotosLibraryClient(authenticatedClient *http.Client) (*PhotosLibraryMediaItemsRepository, error) {
	return NewPhotosLibraryClientWithURL(authenticatedClient, "")
}

// NewPhotosLibraryClientWithURL returns a Repository using PhotosLibrary service with a custom URL.
func NewPhotosLibraryClientWithURL(authenticatedClient *http.Client, url string) (*PhotosLibraryMediaItemsRepository, error) {
	s, err := photoslibrary.New(authenticatedClient)
	if err != nil {
		return nil, err
	}
	if url != "" {
		s.BasePath = url
	}
	return &PhotosLibraryMediaItemsRepository{
		service:  photoslibrary.NewMediaItemsService(s),
		basePath: s.BasePath,
	}, nil
}

// URL returns the repository url.
func (r PhotosLibraryMediaItemsRepository) URL() string {
	return r.basePath
}

// CreateMany creates one or more media items in the repository.
// By default, the media item(s) will be added to the end of the library.
func (r PhotosLibraryMediaItemsRepository) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return r.CreateManyToAlbum(ctx, "", mediaItems)
}

// CreateManyToAlbum creates one or more media item(s) in the repository.
// If an album id is specified, the media item(s) are also added to the album.
// By default, the media item(s) will be added to the end of the library or album.
func (r PhotosLibraryMediaItemsRepository) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
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
	result, err := r.service.BatchCreate(req).Context(ctx).Do()
	if err != nil {
		return []MediaItem{}, err
	}
	mediaItemsResult := make([]MediaItem, len(result.NewMediaItemResults))
	for i, res := range result.NewMediaItemResults {
		// #54: MediaItem is populated if no errors occurred and the media item was created successfully.
		// If an error occurs res.Status should have more data about the error.
		// @see: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchCreate#NewMediaItemResult
		if res.MediaItem != nil {
			mediaItemsResult[i] = toMediaItem(res.MediaItem)
		}
	}
	return mediaItemsResult, nil
}

// Get returns the media item specified based on a given media item id.
func (r PhotosLibraryMediaItemsRepository) Get(ctx context.Context, mediaItemId string) (*MediaItem, error) {
	result, err := r.service.Get(mediaItemId).Context(ctx).Do()
	if err != nil && err.(*googleapi.Error).Code == http.StatusNotFound {
		return &MediaItem{}, ErrNotFound
	}
	if err != nil {
		return &MediaItem{}, ErrServerFailed
	}
	m := toMediaItem(result)
	return &m, nil
}

// ListByAlbum list all media items in the specified album.
func (r PhotosLibraryMediaItemsRepository) ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error) {
	req := &photoslibrary.SearchMediaItemsRequest{
		AlbumId: albumId,
	}

	photosMediaItems := make([]*photoslibrary.MediaItem, 0)
	appendResultsFn := func(result *photoslibrary.SearchMediaItemsResponse) error {
		photosMediaItems = append(photosMediaItems, result.MediaItems...)
		return nil
	}

	if err := r.service.Search(req).Pages(ctx, appendResultsFn); err != nil {
		return []MediaItem{}, err
	}

	mediaItems := make([]MediaItem, len(photosMediaItems))
	for i, item := range photosMediaItems {
		mediaItems[i] = toMediaItem(item)
	}

	return mediaItems, nil
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
