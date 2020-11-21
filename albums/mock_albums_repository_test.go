package albums

import (
	"context"

	duffpl "github.com/duffpl/google-photos-api-client/albums"
)

// MockedAlbumsRepository mocks the an albums repository.
type MockedAlbumsRepository struct {
	BatchAddMediaItemsAllFn    func(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchRemoveMediaItemsAllFn func(albumId string, mediaItemIds []string, ctx context.Context) error
	CreateFn                   func(title string, ctx context.Context) (*duffpl.Album, error)
	GetFn                      func(id string, ctx context.Context) (*duffpl.Album, error)
	ListAllFn                  func(options *duffpl.AlbumsListOptions, ctx context.Context) ([]duffpl.Album, error)
	ListAllAsyncFn             func(options *duffpl.AlbumsListOptions, ctx context.Context) (<-chan duffpl.Album, <-chan error)
}

// BatchAddMediaItemsAll invokes the mock implementation.
func (s MockedAlbumsRepository) BatchAddMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error {
	return s.BatchAddMediaItemsAllFn(albumId, mediaItemIds, ctx)
}

// BatchRemoveMediaItemsAll invokes the mock implementation.
func (s MockedAlbumsRepository) BatchRemoveMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error {
	return s.BatchRemoveMediaItemsAllFn(albumId, mediaItemIds, ctx)
}

// Create invokes the mock implementation.
func (s MockedAlbumsRepository) Create(title string, ctx context.Context) (*duffpl.Album, error) {
	return s.CreateFn(title, ctx)
}

// Get invokes the mock implementation.
func (s MockedAlbumsRepository) Get(id string, ctx context.Context) (*duffpl.Album, error) {
	return s.GetFn(id, ctx)
}

// ListAll invokes the mock implementation.
func (s MockedAlbumsRepository) ListAll(options *duffpl.AlbumsListOptions, ctx context.Context) ([]duffpl.Album, error) {
	return s.ListAllFn(options, ctx)
}

// ListAllAsync invokes the mock implementation.
func (s MockedAlbumsRepository) ListAllAsync(options *duffpl.AlbumsListOptions, ctx context.Context) (<-chan duffpl.Album, <-chan error) {
	return s.ListAllAsyncFn(options, ctx)
}
