package albums

import (
	"context"
	"net/http"

	duffpl "github.com/duffpl/google-photos-api-client/albums"
)

var (
	// Excludes non app created albums. Google Photos doesn't allow manage non created albums through the API.
	// https://developers.google.com/photos/library/guides/manage-albums#adding-items-to-album
	excludeNonAppCreatedData = &duffpl.AlbumsListOptions{ExcludeNonAppCreatedData: true}
)

// DuffplAlbumsClient represents a Albums client using `duffpl/google-photos-api-client`.
type DuffplAlbumsClient interface {
	BatchAddMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchRemoveMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error
	Create(title string, ctx context.Context) (*duffpl.Album, error)
	Get(albumId string, ctx context.Context) (*duffpl.Album, error)
	ListAll(options *duffpl.AlbumsListOptions, ctx context.Context) ([]duffpl.Album, error)
	ListAllAsync(options *duffpl.AlbumsListOptions, ctx context.Context) (<-chan duffpl.Album, <-chan error)
}

type AlbumRepository struct {
	duffplAlbumsClient DuffplAlbumsClient
}

func (ar AlbumRepository) AddManyItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return ar.duffplAlbumsClient.BatchAddMediaItemsAll(albumId, mediaItemIds, ctx)
}

func (ar AlbumRepository) RemoveManyItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return ar.duffplAlbumsClient.BatchRemoveMediaItemsAll(albumId, mediaItemIds, ctx)
}

func (ar AlbumRepository) Create(ctx context.Context, title string) (*Album, error) {
	a, err := ar.duffplAlbumsClient.Create(title, ctx)
	if err != nil {
		return &NullAlbum, err
	}
	album := ar.convertDuffplAlbumToAlbum(a)
	return &album, nil
}

func (ar AlbumRepository) Get(ctx context.Context, albumId string) (*Album, error) {
	a, err := ar.duffplAlbumsClient.Get(albumId, ctx)
	if err != nil {
		return &NullAlbum, err
	}
	album := ar.convertDuffplAlbumToAlbum(a)
	return &album, nil
}

func (ar AlbumRepository) GetByTitle(ctx context.Context, title string) (*Album, error) {
	albumsC, errorsC := ar.duffplAlbumsClient.ListAllAsync(excludeNonAppCreatedData, ctx)
	for {
		select {
		case item, ok := <-albumsC:
			if !ok {
				return &NullAlbum, ErrAlbumNotFound // there aren't more albums, album not found
			}
			if item.Title == title {
				album := ar.convertDuffplAlbumToAlbum(&item)
				return &album, nil // found, cache it and return it
			}
		case err := <-errorsC:
			return &NullAlbum, err
		}
	}
}

func (ar AlbumRepository) ListAll(ctx context.Context) ([]Album, error) {
	albums := make([]Album, 0)
	result, err := ar.duffplAlbumsClient.ListAll(excludeNonAppCreatedData, ctx)
	if err != nil {
		return albums, err
	}
	for _, a := range result {
		albums = append(albums, ar.convertDuffplAlbumToAlbum(&a))
	}
	return albums, nil
}


// NewDuffplAlbumRepository implements AlbumRepository using https://github.com/duffpl/google-photos-api-client library.
func NewDuffplAlbumRepository(authenticatedClient *http.Client) AlbumRepository {
	return AlbumRepository{
		duffplAlbumsClient: duffpl.NewHttpAlbumsService(authenticatedClient),
	}
}

func (ar AlbumRepository) convertDuffplAlbumToAlbum(a *duffpl.Album) Album {
	return Album{
		ID:                    a.ID,
		Title:                 a.Title,
		ProductURL:            a.ProductURL,
		IsWriteable:           a.IsWriteable,
		MediaItemsCount:       a.MediaItemsCount,
		CoverPhotoBaseURL:     a.CoverPhotoBaseURL,
		CoverPhotoMediaItemID: a.CoverPhotoMediaItemID,
	}
}
