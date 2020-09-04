package mock

import (
	"context"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

type PhotoService struct {
	ListAlbumsFn      func(ctx context.Context, pageSize int64, pageToken string) (*photoslibrary.ListAlbumsResponse, error)
	ListAlbumsInvoked bool

	CreateAlbumFn      func(ctx context.Context, request *photoslibrary.CreateAlbumRequest) (*photoslibrary.Album, error)
	CreateAlbumInvoked bool

	CreateMediaItemsFn      func(ctx context.Context, request *photoslibrary.BatchCreateMediaItemsRequest) (*photoslibrary.BatchCreateMediaItemsResponse, error)
	CreateMediaItemsInvoked bool
}

// ListAlbums invokes the mock implementation and marks the function as invoked.
func (s *PhotoService) ListAlbums(ctx context.Context, pageSize int64, pageToken string) (*photoslibrary.ListAlbumsResponse, error) {
	s.ListAlbumsInvoked = true
	return s.ListAlbumsFn(ctx, pageSize, pageToken)
}

// CreateAlbum invokes the mock implementation and marks the function as invoked.
func (s *PhotoService) CreateAlbum(ctx context.Context, request *photoslibrary.CreateAlbumRequest) (*photoslibrary.Album, error) {
	s.CreateAlbumInvoked = true
	return s.CreateAlbumFn(ctx, request)
}

func (s *PhotoService) CreateMediaItems(ctx context.Context, request *photoslibrary.BatchCreateMediaItemsRequest) (*photoslibrary.BatchCreateMediaItemsResponse, error) {
	s.CreateMediaItemsInvoked = true
	return s.CreateMediaItemsFn(ctx, request)
}
