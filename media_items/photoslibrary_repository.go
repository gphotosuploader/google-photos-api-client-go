package media_items

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// PhotosLibraryClient represents a media items service using `gphotosuploader/googlemirror/api/photoslibrary`.
type PhotosLibraryClient interface {
	BatchCreate(batchcreatemediaitemsrequest *photoslibrary.BatchCreateMediaItemsRequest) *photoslibrary.MediaItemsBatchCreateCall
	Get(mediaItemId string) *photoslibrary.MediaItemsGetCall
	Search(searchmediaitemsrequest *photoslibrary.SearchMediaItemsRequest) *photoslibrary.MediaItemsSearchCall
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

// URL returns the media items repository url.
func (r PhotosLibraryMediaItemsRepository) URL() string {
	return r.basePath
}

// CreateMany creates one or more media items in the repository.
// By default the media item(s) will be added to the end of the library.
func (r PhotosLibraryMediaItemsRepository) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return r.CreateManyToAlbum(ctx, "", mediaItems)
}

// CreateManyToAlbum creates one or more media item(s) in the repository.
// If an album id is specified, the media item(s) are also added to the album.
// By default the media item(s) will be added to the end of the library or album.
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
		m := res.MediaItem
		mediaItemsResult[i] = r.convertPhotosLibraryMediaItemToMediaItem(m)
	}
	return mediaItemsResult, nil
}

// Get returns the media item specified based on a given media item id.
func (r PhotosLibraryMediaItemsRepository) Get(ctx context.Context, mediaItemId string) (*MediaItem, error) {
	result, err := r.service.Get(mediaItemId).Context(ctx).Do()
	if err != nil {
		return &MediaItem{}, err
	}
	m := r.convertPhotosLibraryMediaItemToMediaItem(result)
	return &m, nil
}

// ListByAlbum list all media items in the specified album.
func (r PhotosLibraryMediaItemsRepository) ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error) {
	req := &photoslibrary.SearchMediaItemsRequest{
		AlbumId: albumId,
	}
	result, err := r.service.Search(req).Context(ctx).Do()
	if err != nil {
		return []MediaItem{}, err
	}
	mediaItemsResult := make([]MediaItem, 0)
	for _, res := range result.MediaItems {
		mediaItemsResult = append(mediaItemsResult,
			r.convertPhotosLibraryMediaItemToMediaItem(res))
	}
	return mediaItemsResult, nil
}

func (r PhotosLibraryMediaItemsRepository) convertPhotosLibraryMediaItemToMediaItem(m *photoslibrary.MediaItem) MediaItem {
	return MediaItem{
		ID:         m.Id,
		ProductURL: m.ProductUrl,
		BaseURL:    m.BaseUrl,
		MimeType:   m.MimeType,
		MediaMetadata: MediaMetadata{
			CreationTime: m.MediaMetadata.CreationTime,
			Width:        strconv.FormatInt(m.MediaMetadata.Width, 10),
			Height:       strconv.FormatInt(m.MediaMetadata.Height, 10),
		},
		Filename: m.Filename,
	}
}
