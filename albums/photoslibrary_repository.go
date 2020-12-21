package albums

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// PhotosLibraryClient represents a AlbumsService using `gphotosuploader/googlemirror/api/photoslibrary`.
type PhotosLibraryClient interface {
	BatchAddMediaItems(albumId string, albumbatchaddmediaitemsrequest *photoslibrary.AlbumBatchAddMediaItemsRequest) *photoslibrary.AlbumBatchAddMediaItemsCall
	Create(createalbumrequest *photoslibrary.CreateAlbumRequest) *photoslibrary.AlbumsCreateCall
	Get(albumId string) *photoslibrary.AlbumsGetCall
	List() *photoslibrary.AlbumsListCall
}

type PhotosLibraryAlbumsRepository struct {
	service  PhotosLibraryClient
	basePath string
}

// NewPhotosLibraryClient returns a Repository using PhotosLibrary service.
func NewPhotosLibraryClient(authenticatedClient *http.Client) (*PhotosLibraryAlbumsRepository, error) {
	return NewPhotosLibraryClientWithURL(authenticatedClient, "")
}

// NewPhotosLibraryClientWithURL returns a Repository using PhotosLibrary service with a custom URL.
func NewPhotosLibraryClientWithURL(authenticatedClient *http.Client, url string) (*PhotosLibraryAlbumsRepository, error) {
	s, err := photoslibrary.New(authenticatedClient)
	if err != nil {
		return nil, err
	}
	if url != "" {
		s.BasePath = url
	}
	return &PhotosLibraryAlbumsRepository{
		service:  photoslibrary.NewAlbumsService(s),
		basePath: s.BasePath,
	}, nil
}

func (ar PhotosLibraryAlbumsRepository) URL() string {
	return ar.basePath
}

func (ar PhotosLibraryAlbumsRepository) AddManyItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	req := &photoslibrary.AlbumBatchAddMediaItemsRequest{
		MediaItemIds: mediaItemIds,
	}
	_, err := ar.service.BatchAddMediaItems(albumId, req).Context(ctx).Do()
	return err
}

func (ar PhotosLibraryAlbumsRepository) RemoveManyItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	panic("not implemented on google mirror library")
}

func (ar PhotosLibraryAlbumsRepository) Create(ctx context.Context, title string) (*Album, error) {
	req := &photoslibrary.CreateAlbumRequest{
		Album: &photoslibrary.Album{Title: title},
	}
	res, err := ar.service.Create(req).Context(ctx).Do()
	if err != nil {
		return &NullAlbum, err
	}
	album := ar.convertPhotosLibraryAlbumToAlbum(res)
	return &album, nil
}

func (ar PhotosLibraryAlbumsRepository) Get(ctx context.Context, albumId string) (*Album, error) {
	res, err := ar.service.Get(albumId).Context(ctx).Do()
	if err != nil {
		return &NullAlbum, ErrAlbumNotFound
	}
	album := ar.convertPhotosLibraryAlbumToAlbum(res)
	return &album, nil
}

func (ar PhotosLibraryAlbumsRepository) ListAll(ctx context.Context) ([]Album, error) {
	albumsResult := make([]Album, 0)
	err := ar.service.List().ExcludeNonAppCreatedData().Pages(ctx, func(response *photoslibrary.ListAlbumsResponse) error {
		for _, res := range response.Albums {
			albumsResult = append(albumsResult, ar.convertPhotosLibraryAlbumToAlbum(res))
		}
		return nil
	})
	return albumsResult, err
}

func (ar PhotosLibraryAlbumsRepository) GetByTitle(ctx context.Context, title string) (*Album, error) {
	ErrAlbumWasFound := fmt.Errorf("album was found")
	var albumResult Album
	if err := ar.service.List().ExcludeNonAppCreatedData().Pages(ctx, func(response *photoslibrary.ListAlbumsResponse) error {
		for _, res := range response.Albums {
			if res.Title == title {
				albumResult = ar.convertPhotosLibraryAlbumToAlbum(res)
				return ErrAlbumWasFound
			}
		}
		return nil
	}); err == ErrAlbumWasFound {
		return &albumResult, nil
	}
	return &NullAlbum, ErrAlbumNotFound
}

func (ar PhotosLibraryAlbumsRepository) convertPhotosLibraryAlbumToAlbum(a *photoslibrary.Album) Album {
	return Album{
		ID:                    a.Id,
		Title:                 a.Title,
		ProductURL:            a.ProductUrl,
		IsWriteable:           a.IsWriteable,
		MediaItemsCount:       strconv.FormatInt(a.TotalMediaItems, 10),
		CoverPhotoBaseURL:     a.CoverPhotoBaseUrl,
		CoverPhotoMediaItemID: "",
	}
}
