package mocks

import (
	"context"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
)

// MockedAlbumsService mocks the albums service.
type MockedAlbumsService struct {
	AddMediaItemsFn    func(ctx context.Context, albumId string, mediaItemIds []string) error
	RemoveMediaItemsFn func(ctx context.Context, albumId string, mediaItemIds []string) error
	CreateFn           func(ctx context.Context, title string) (*albums.Album, error)
	GetByIdFn          func(ctx context.Context, id string) (*albums.Album, error)
	GetByTitleFn       func(ctx context.Context, title string) (*albums.Album, error)
	ListFn             func(ctx context.Context) ([]albums.Album, error)
}

// AddMediaItems invokes the mock implementation.
func (m MockedAlbumsService) AddMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return m.AddMediaItemsFn(ctx, albumId, mediaItemIds)
}

// RemoveMediaItems invokes the mock implementation.
func (m MockedAlbumsService) RemoveMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return m.RemoveMediaItemsFn(ctx, albumId, mediaItemIds)
}

// Create invokes the mock implementation.
func (m MockedAlbumsService) Create(ctx context.Context, title string) (*albums.Album, error) {
	return m.CreateFn(ctx, title)
}

// GetById invokes the mock implementation.
func (m MockedAlbumsService) GetById(ctx context.Context, id string) (*albums.Album, error) {
	return m.GetByIdFn(ctx, id)
}

// GetByTitle invokes the mock implementation.
func (m MockedAlbumsService) GetByTitle(ctx context.Context, title string) (*albums.Album, error) {
	return m.GetByTitleFn(ctx, title)
}

// List invokes the mock implementation.
func (m MockedAlbumsService) List(ctx context.Context) ([]albums.Album, error) {
	return m.ListFn(ctx)
}
