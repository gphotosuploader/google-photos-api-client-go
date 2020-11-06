package mock

import (
	"context"

	"github.com/duffpl/google-photos-api-client/albums"
)

// AlbumsRepository mocks the an albums repository.
type AlbumsRepository struct {
	BatchAddMediaItemsAllFn      func(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchAddMediaItemsAllInvoked bool

	BatchRemoveMediaItemsAllFn      func(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchRemoveMediaItemsAllInvoked bool

	CreateFn      func(title string, ctx context.Context) (*albums.Album, error)
	CreateInvoked bool

	GetFn      func(id string, ctx context.Context) (*albums.Album, error)
	GetInvoked bool

	ListAllFn      func(options *albums.AlbumsListOptions, ctx context.Context) ([]albums.Album, error)
	ListAllInvoked bool

	ListAllAsyncFn      func(options *albums.AlbumsListOptions, ctx context.Context) (<-chan albums.Album, <-chan error)
	ListAllAsyncInvoked bool

	PatchFn      func(album albums.Album, updateMask []albums.Field, ctx context.Context) (*albums.Album, error)
	PatchInvoked bool
}

// BatchAddMediaItemsAll invokes the mock implementation and marks the function as invoked.
func (s AlbumsRepository) BatchAddMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error {
	s.BatchAddMediaItemsAllInvoked = true
	return s.BatchAddMediaItemsAllFn(albumId, mediaItemIds, ctx)
}

// BatchRemoveMediaItemsAll invokes the mock implementation and marks the function as invoked.
func (s AlbumsRepository) BatchRemoveMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error {
	s.BatchRemoveMediaItemsAllInvoked = true
	return s.BatchRemoveMediaItemsAllFn(albumId, mediaItemIds, ctx)
}

// Create invokes the mock implementation and marks the function as invoked.
func (s AlbumsRepository) Create(title string, ctx context.Context) (*albums.Album, error) {
	s.CreateInvoked = true
	return s.CreateFn(title, ctx)
}

// Get invokes the mock implementation and marks the function as invoked.
func (s AlbumsRepository) Get(id string, ctx context.Context) (*albums.Album, error) {
	s.GetInvoked = true
	return s.GetFn(id, ctx)
}

// ListAll invokes the mock implementation and marks the function as invoked.
func (s AlbumsRepository) ListAll(options *albums.AlbumsListOptions, ctx context.Context) ([]albums.Album, error) {
	s.ListAllInvoked = true
	return s.ListAllFn(options, ctx)
}

// ListAllAsync invokes the mock implementation and marks the function as invoked.
func (s AlbumsRepository) ListAllAsync(options *albums.AlbumsListOptions, ctx context.Context) (<-chan albums.Album, <-chan error) {
	s.ListAllAsyncInvoked = true
	return s.ListAllAsyncFn(options, ctx)
}

// Patch invokes the mock implementation and marks the function as invoked.
func (s AlbumsRepository) Patch(album albums.Album, updateMask []albums.Field, ctx context.Context) (*albums.Album, error) {
	s.PatchInvoked = true
	return s.PatchFn(album, updateMask, ctx)
}
