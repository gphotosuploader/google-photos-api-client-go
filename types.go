package gphotos

import (
	"context"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
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
