package media_items

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// PhotosLibraryClient represents a MediaItemsService using `gphotosuploader/googlemirror/api/photoslibrary`.
type PhotosLibraryClient interface {
	BatchCreate(batchcreatemediaitemsrequest *photoslibrary.BatchCreateMediaItemsRequest) *photoslibrary.MediaItemsBatchCreateCall
	Get(mediaItemId string) *photoslibrary.MediaItemsGetCall
	Search(searchmediaitemsrequest *photoslibrary.SearchMediaItemsRequest) *photoslibrary.MediaItemsSearchCall
}

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

func (r PhotosLibraryMediaItemsRepository) URL() string {
	return r.basePath
}

func (r PhotosLibraryMediaItemsRepository) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return r.CreateManyToAlbum(ctx, "", mediaItems)
}

func (r PhotosLibraryMediaItemsRepository) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	newMediaItems := make([]*photoslibrary.NewMediaItem, 0)
	for _, mediaItem := range mediaItems {
		newMediaItems = append(newMediaItems, &photoslibrary.NewMediaItem{
			SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: mediaItem.UploadToken},
		})
	}
	req := &photoslibrary.BatchCreateMediaItemsRequest{
		AlbumId:       albumId,
		NewMediaItems: newMediaItems,
	}
	result, err := r.service.BatchCreate(req).Context(ctx).Do()
	if err != nil {
		return []MediaItem{}, err
	}
	mediaItemsResult := make([]MediaItem, 0)
	for _, res := range result.NewMediaItemResults {
		m := res.MediaItem
		mediaItemsResult = append(mediaItemsResult, r.convertPhotosLibraryMediaItemToMediaItem(m))
	}
	return mediaItemsResult, nil
}

func (r PhotosLibraryMediaItemsRepository) Get(ctx context.Context, mediaItemId string) (*MediaItem, error) {
	result, err := r.service.Get(mediaItemId).Context(ctx).Do()
	if err != nil {
		return &MediaItem{}, err
	}
	m := r.convertPhotosLibraryMediaItemToMediaItem(result)
	return &m, nil
}

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
