package gphotos

import (
	"context"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/media_items"
)

// OAuth2 scopes used by this API.
const (
	// PhotoslibraryScope View and manage your Google Photos library
	PhotoslibraryScope = "https://www.googleapis.com/auth/photoslibrary"

	// PhotoslibraryAppendonlyScope Add to your Google Photos library
	PhotoslibraryAppendonlyScope = "https://www.googleapis.com/auth/photoslibrary.appendonly"

	// PhotoslibraryReadonlyScope View your Google Photos library
	PhotoslibraryReadonlyScope = "https://www.googleapis.com/auth/photoslibrary.readonly"

	// PhotoslibraryReadonlyAppcreateddataScope Manage photos added by this app
	PhotoslibraryReadonlyAppcreateddataScope = "https://www.googleapis.com/auth/photoslibrary.readonly.appcreateddata"

	// PhotoslibrarySharingScope Manage and add to shared albums on your behalf
	PhotoslibrarySharingScope = "https://www.googleapis.com/auth/photoslibrary.sharing"
)

// AlbumsService represents a Google Photos client for albums management.
type AlbumsService interface {
	AddMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error
	RemoveMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error
	Create(ctx context.Context, title string) (*albums.Album, error)
	GetById(ctx context.Context, id string) (*albums.Album, error)
	GetByTitle(ctx context.Context, title string) (*albums.Album, error)
	List(ctx context.Context) ([]albums.Album, error)
}

// MediaItemsService represents a Google Photos client for media management.
type MediaItemsService interface {
	Create(ctx context.Context, mediaItem media_items.SimpleMediaItem) (media_items.MediaItem, error)
	CreateMany(ctx context.Context, mediaItems []media_items.SimpleMediaItem) ([]media_items.MediaItem, error)
	CreateToAlbum(ctx context.Context, albumId string, mediaItem media_items.SimpleMediaItem) (media_items.MediaItem, error)
	CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []media_items.SimpleMediaItem) ([]media_items.MediaItem, error)
	Get(ctx context.Context, mediaItemId string) (*media_items.MediaItem, error)
	ListByAlbum(ctx context.Context, albumId string) ([]media_items.MediaItem, error)
}

// MediaUploader represents a Google Photos client fo media upload.
type MediaUploader interface {
	UploadFile(ctx context.Context, filePath string) (uploadToken string, err error)
}
