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

// GooglePhotosRepository implements Repository using `gphotosuploader/googlemirror/api/photoslibrary`.
type GooglePhotosRepository struct {
	gPhotosClient PhotosLibraryClient
}

// NewPhotosLibraryClient returns a Repository using PhotosLibrary service.
func NewPhotosLibraryClient(authenticatedClient *http.Client) (*GooglePhotosRepository, error) {
	service, err := photoslibrary.New(authenticatedClient)
	if err != nil {
		return nil, err
	}
	return &GooglePhotosRepository{
		gPhotosClient: photoslibrary.NewMediaItemsService(service),
	}, nil
}

func (r GooglePhotosRepository) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return r.CreateManyToAlbum(ctx, "", mediaItems)
}

func (r GooglePhotosRepository) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
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
	result, err := r.gPhotosClient.BatchCreate(req).Context(ctx).Do()
	if err != nil {
		return []MediaItem{}, err
	}
	mediaItemsResult := make([]MediaItem, 0)
	for _, res := range result.NewMediaItemResults {
		mediaItemsResult = append(mediaItemsResult, MediaItem{
			ID:         res.MediaItem.Id,
			ProductURL: res.MediaItem.ProductUrl,
			BaseURL:    res.MediaItem.BaseUrl,
			MimeType:   res.MediaItem.MimeType,
			MediaMetadata: MediaMetadata{
				CreationTime: res.MediaItem.MediaMetadata.CreationTime,
				Width:        strconv.FormatInt(res.MediaItem.MediaMetadata.Width, 10),
				Height:       strconv.FormatInt(res.MediaItem.MediaMetadata.Height, 10),
			},
			Filename: res.MediaItem.Filename,
		})
	}
	return mediaItemsResult, nil
}

func (r GooglePhotosRepository) Get(ctx context.Context, mediaItemId string) (*MediaItem, error) {
	result, err := r.gPhotosClient.Get(mediaItemId).Context(ctx).Do()
	if err != nil {
		return &MediaItem{}, err
	}
	return &MediaItem{
		ID:          result.Id,
		Description: result.Description,
		ProductURL:  result.ProductUrl,
		BaseURL:     result.BaseUrl,
		MimeType:    result.MimeType,
		MediaMetadata: MediaMetadata{
			CreationTime: result.MediaMetadata.CreationTime,
			Width:        strconv.FormatInt(result.MediaMetadata.Width, 10),
			Height:       strconv.FormatInt(result.MediaMetadata.Height, 10),
		},
		Filename: "",
	}, nil
}

func (r GooglePhotosRepository) ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error) {
	req := &photoslibrary.SearchMediaItemsRequest{
		AlbumId: albumId,
	}
	result, err := r.gPhotosClient.Search(req).Context(ctx).Do()
	if err != nil {
		return []MediaItem{}, err
	}
	mediaItemsResult := make([]MediaItem, 0)
	for _, res := range result.MediaItems {
		mediaItemsResult = append(mediaItemsResult, MediaItem{
			ID:         res.Id,
			ProductURL: res.ProductUrl,
			BaseURL:    res.BaseUrl,
			MimeType:   res.MimeType,
			MediaMetadata: MediaMetadata{
				CreationTime: res.MediaMetadata.CreationTime,
				Width:        strconv.FormatInt(res.MediaMetadata.Width, 10),
				Height:       strconv.FormatInt(res.MediaMetadata.Height, 10),
			},
			Filename: res.Filename,
		})
	}
	return mediaItemsResult, nil
}
