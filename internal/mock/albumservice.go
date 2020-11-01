package mock

import (
	"context"

	"github.com/duffpl/google-photos-api-client/albums"
)

// AlbumService mocks a album service.
type AlbumService struct {
	AddEnrichmentFn      func(albumId string, enrichment albums.NewEnrichmentItem, ctx context.Context) (*albums.EnrichmentItem, error)
	AddEnrichmentInvoked bool

	BatchAddMediaItemsFn      func(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchAddMediaItemsInvoked bool

	BatchAddMediaItemsAllFn      func(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchAddMediaItemsAllInvoked bool

	BatchRemoveMediaItemsFn      func(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchRemoveMediaItemsInvoked bool

	BatchRemoveMediaItemsAllFn      func(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchRemoveMediaItemsAllInvoked bool

	CreateFn      func(title string, ctx context.Context) (*albums.Album, error)
	CreateInvoked bool

	GetFn      func(id string, ctx context.Context) (*albums.Album, error)
	GetInvoked bool

	ListFn      func(options *albums.AlbumsListOptions, pageToken string, ctx context.Context) (result []albums.Album, nextPageToken string, err error)
	ListInvoked bool

	ListAllFn      func(options *albums.AlbumsListOptions, ctx context.Context) ([]albums.Album, error)
	ListAllInvoked bool

	ListAllAsyncFn      func(options *albums.AlbumsListOptions, ctx context.Context) (<-chan albums.Album, <-chan error)
	ListAllAsyncInvoked bool

	PatchFn      func(album albums.Album, updateMask []albums.Field, ctx context.Context) (*albums.Album, error)
	PatchInvoked bool

	ShareFn      func(id string, options albums.SharedAlbumOptions, ctx context.Context) (*albums.AlbumShareInfo, error)
	ShareInvoked bool

	UnshareFn      func(id string, ctx context.Context) error
	UnshareInvoked bool
}

// AddEnrichment invokes the mock implementation and marks the function as invoked.
func (s AlbumService) AddEnrichment(albumId string, enrichment albums.NewEnrichmentItem, ctx context.Context) (*albums.EnrichmentItem, error) {
	s.AddEnrichmentInvoked = true
	return s.AddEnrichmentFn(albumId, enrichment, ctx)
}

// BatchAddMediaItems invokes the mock implementation and marks the function as invoked.
func (s AlbumService) BatchAddMediaItems(albumId string, mediaItemIds []string, ctx context.Context) error {
	s.BatchAddMediaItemsInvoked = true
	return s.BatchAddMediaItemsFn(albumId, mediaItemIds, ctx)
}

// BatchAddMediaItemsAll invokes the mock implementation and marks the function as invoked.
func (s AlbumService) BatchAddMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error {
	s.BatchAddMediaItemsAllInvoked = true
	return s.BatchAddMediaItemsAllFn(albumId, mediaItemIds, ctx)
}

// BatchRemoveMediaItems invokes the mock implementation and marks the function as invoked.
func (s AlbumService) BatchRemoveMediaItems(albumId string, mediaItemIds []string, ctx context.Context) error {
	s.BatchRemoveMediaItemsInvoked = true
	return s.BatchRemoveMediaItemsFn(albumId, mediaItemIds, ctx)
}

// BatchRemoveMediaItemsAll invokes the mock implementation and marks the function as invoked.
func (s AlbumService) BatchRemoveMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error {
	s.BatchRemoveMediaItemsAllInvoked = true
	return s.BatchRemoveMediaItemsAllFn(albumId, mediaItemIds, ctx)
}

// Create invokes the mock implementation and marks the function as invoked.
func (s AlbumService) Create(title string, ctx context.Context) (*albums.Album, error) {
	s.CreateInvoked = true
	return s.CreateFn(title, ctx)
}

// Get invokes the mock implementation and marks the function as invoked.
func (s AlbumService) Get(id string, ctx context.Context) (*albums.Album, error) {
	s.GetInvoked = true
	return s.GetFn(id, ctx)
}

// List invokes the mock implementation and marks the function as invoked.
func (s AlbumService) List(options *albums.AlbumsListOptions, pageToken string, ctx context.Context) (result []albums.Album, nextPageToken string, err error) {
	s.ListInvoked = true
	return s.ListFn(options, pageToken, ctx)
}

// ListAll invokes the mock implementation and marks the function as invoked.
func (s AlbumService) ListAll(options *albums.AlbumsListOptions, ctx context.Context) ([]albums.Album, error) {
	s.ListAllInvoked = true
	return s.ListAllFn(options, ctx)
}

// ListAllAsync invokes the mock implementation and marks the function as invoked.
func (s AlbumService) ListAllAsync(options *albums.AlbumsListOptions, ctx context.Context) (<-chan albums.Album, <-chan error) {
	s.ListAllAsyncInvoked = true
	return s.ListAllAsyncFn(options, ctx)
}

// Patch invokes the mock implementation and marks the function as invoked.
func (s AlbumService) Patch(album albums.Album, updateMask []albums.Field, ctx context.Context) (*albums.Album, error) {
	s.PatchInvoked = true
	return s.PatchFn(album, updateMask, ctx)
}

// Share invokes the mock implementation and marks the function as invoked.
func (s AlbumService) Share(id string, options albums.SharedAlbumOptions, ctx context.Context) (*albums.AlbumShareInfo, error) {
	s.ShareInvoked = true
	return s.ShareFn(id, options, ctx)
}

// Unshare invokes the mock implementation and marks the function as invoked.
func (s AlbumService) Unshare(id string, ctx context.Context) error {
	s.UnshareInvoked = true
	return s.UnshareFn(id, ctx)
}
