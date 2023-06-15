package albums

import (
	"context"
	"errors"
	"fmt"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"net/http"
	"strconv"
)

// Config holds the configuration parameters for the service.
type Config struct {
	// Client should have all oAuth credentials in place.
	Client *http.Client
	URL    string
}

// Service implements an albums Google Photos client.
type Service struct {
	photos PhotosLibraryClient
}

// PhotosLibraryClient represents a Google Photos client using `gphotosuploader/googlemirror/api/photoslibrary`.
type PhotosLibraryClient interface {
	BatchAddMediaItems(albumId string, albumBatchAddMediaItemsRequest *photoslibrary.AlbumBatchAddMediaItemsRequest) *photoslibrary.AlbumBatchAddMediaItemsCall
	Create(createAlbumRequest *photoslibrary.CreateAlbumRequest) *photoslibrary.AlbumsCreateCall
	Get(albumId string) *photoslibrary.AlbumsGetCall
	List() *photoslibrary.AlbumsListCall
}

var (
	// NullAlbum is a zero value Album.
	NullAlbum = Album{}

	// ErrAlbumNotFound is the error returned when an album is not found.
	ErrAlbumNotFound = errors.New("album not found")
)

// AddMediaItems adds multiple media item(s) to the specified album.
func (s Service) AddMediaItems(ctx context.Context, albumID string, mediaItemIDs []string) error {
	req := &photoslibrary.AlbumBatchAddMediaItemsRequest{
		MediaItemIds: mediaItemIDs,
	}
	_, err := s.photos.BatchAddMediaItems(albumID, req).Context(ctx).Do()
	return err
}

// RemoveMediaItems removes multiple media item(s) from the specified album.
func (s Service) RemoveMediaItems(ctx context.Context, albumID string, mediaItemIDs []string) error {
	panic("not implemented on google mirror library")
}

// Create adds a new album to the repo.
func (s Service) Create(ctx context.Context, title string) (*Album, error) {
	req := &photoslibrary.CreateAlbumRequest{
		Album: &photoslibrary.Album{Title: title},
	}
	res, err := s.photos.Create(req).Context(ctx).Do()
	if err != nil {
		return &NullAlbum, err
	}
	album := toAlbum(res)
	return &album, nil
}

// GetById fetches an album from the repository by id.
// It returns ErrAlbumNotFound if the album does not exist.
func (s Service) GetById(ctx context.Context, albumID string) (*Album, error) {
	res, err := s.photos.Get(albumID).Context(ctx).Do()
	if err != nil {
		return &NullAlbum, fmt.Errorf("%s: %w", albumID, ErrAlbumNotFound)
	}
	album := toAlbum(res)
	return &album, nil
}

// GetByTitle fetches an album from the repository by title.
// Ir returns ErrAlbumNotFound if the album does not exist.
func (s Service) GetByTitle(ctx context.Context, title string) (*Album, error) {
	errAlbumWasFound := errors.New("album was found")
	var result *Album
	if err := s.photos.List().ExcludeNonAppCreatedData().Pages(ctx, func(response *photoslibrary.ListAlbumsResponse) error {
		if album, found := findByTitle(title, response.Albums); found {
			result = album
			return errAlbumWasFound
		}
		return nil
	}); errors.Is(err, errAlbumWasFound) {
		return result, nil
	}
	return &NullAlbum, fmt.Errorf("%s: %w", title, ErrAlbumNotFound)
}

// List fetches all the albums from the repository.
func (s Service) List(ctx context.Context) ([]Album, error) {
	albumsResult := make([]Album, 0)
	albumsListCall := s.photos.List().PageSize(maxItemsPerPage).ExcludeNonAppCreatedData()
	err := albumsListCall.Pages(ctx, func(response *photoslibrary.ListAlbumsResponse) error {
		for _, res := range response.Albums {
			albumsResult = append(albumsResult, toAlbum(res))
		}
		return nil
	})
	return albumsResult, err
}

func NewService(config Config) (*Service, error) {
	client := config.Client

	if client == nil {
		client = http.DefaultClient
	}

	photosClient, err := newPhotosLibraryClient(client, config.URL)
	if err != nil {
		return nil, err
	}

	service := &Service{
		photos: photosClient,
	}

	return service, nil
}

func newPhotosLibraryClient(authenticatedClient *http.Client, url string) (*photoslibrary.AlbumsService, error) {
	s, err := photoslibrary.New(authenticatedClient)
	if err != nil {
		return nil, err
	}
	if url != "" {
		s.BasePath = url
	}
	return s.Albums, nil
}

func findByTitle(title string, albums []*photoslibrary.Album) (*Album, bool) {
	for _, a := range albums {
		if a.Title == title {
			album := toAlbum(a)
			return &album, true
		}
	}
	return &NullAlbum, false
}

func toAlbum(a *photoslibrary.Album) Album {
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
