package gphotos

import (
	"context"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

type Service interface {
	ListAlbums(ctx context.Context, pageSize int64, pageToken string) (*photoslibrary.ListAlbumsResponse, error)
	CreateAlbum(ctx context.Context, request *photoslibrary.CreateAlbumRequest) (*photoslibrary.Album, error)

	CreateMediaItems(ctx context.Context, request *photoslibrary.BatchCreateMediaItemsRequest) (*photoslibrary.BatchCreateMediaItemsResponse, error)
}
