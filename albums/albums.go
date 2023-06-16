package albums

import (
	"context"
	"errors"
	"fmt"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"net/http"
)

// An Album represents a Google Photos album.
// Albums are a container for media items.
//
// See: https://developers.google.com/photos/library/reference/rest/v1/albums.
type Album struct {
	// CoverPhotoBaseURL: [Output only] A BaseURL to the cover photo's bytes.
	// This should not be used as is.
	// Parameters should be appended to this BaseURL before use. For example,
	// '=w2048-h1024' will set the dimensions of the cover photo to have a
	// width of 2048 px and height of 1024 px.
	CoverPhotoBaseURL string

	// Id: [Output only] Identifier for the album. This is a persistent
	// identifier that can be used to identify this album.
	ID string

	// IsWriteable: [Output only] True if media items can be created in the
	// album.
	// This field is based on the scopes granted and permissions of the
	// album. If the scopes are changed or permissions of the album are changed, this
	// field will be updated.
	IsWriteable bool

	// ProductURL: [Output only] Google Photos BaseURL for the album. The user
	// needs to be signed in to their Google Photos account to access this link.
	ProductURL string

	// Title: Name of the album displayed to the user in their Google Photos
	// account.
	// This string should not be more than 500 characters.
	Title string

	// TotalMediaItems: [Output only] The number of media items in the album.
	TotalMediaItems int64
}

// Config holds the configuration parameters for the service.
type Config struct {
	// HTTP client used to communicate with the API.
	Client *http.Client

	// [Optional] Base URL for API requests.
	// BaseURL should always be specified with a trailing slash.
	BaseURL string

	// [Optional] User agent used when communicating with the Google Photos API.
	UserAgent string
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

func New(config Config) (*Service, error) {
	if config.Client == nil {
		return nil, errors.New("client is nil")
	}

	s, err := photoslibrary.New(config.Client)
	if err != nil {
		return nil, err
	}

	if config.BaseURL != "" {
		s.BasePath = config.BaseURL
	}

	if config.UserAgent != "" {
		s.UserAgent = config.UserAgent
	}

	service := &Service{
		photos: s.Albums,
	}

	return service, nil
}

// AddMediaItems add one or more existing media items to an existing Album.
func (s *Service) AddMediaItems(ctx context.Context, albumID string, mediaItemIDs []string) error {

	// TODO: There's a limit of 50 media items per call. Split in multiple calls if more are provided.

	req := &photoslibrary.AlbumBatchAddMediaItemsRequest{
		MediaItemIds: mediaItemIDs,
	}
	_, err := s.photos.BatchAddMediaItems(albumID, req).Context(ctx).Do()
	return err
}

// Create creates an album in Google Photos.
func (s *Service) Create(ctx context.Context, title string) (*Album, error) {
	req := &photoslibrary.CreateAlbumRequest{
		Album: &photoslibrary.Album{Title: title},
	}
	res, err := s.photos.Create(req).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	album := toAlbum(res)
	return &album, nil
}

// GetById returns the album specified by the given album id.
//
// Returns [ErrAlbumNotFound] if the album does not exist.
func (s *Service) GetById(ctx context.Context, albumID string) (*Album, error) {
	res, err := s.photos.Get(albumID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("getById %s: %w", albumID, ErrAlbumNotFound)
	}
	album := toAlbum(res)
	return &album, nil
}

// maxAlbumsPerPage is the maximum number of albums per pages.
// Fewer albums might be returned than the specified number.
// See https://developers.google.com/photos/library/guides/list#pagination.
const maxAlbumsPerPage = 50

// GetByTitle look for an album with the specified album id into the list of all albums.
// It lists paginates all albums until finding one with the matching title.
//
// Returns [ErrAlbumNotFound] if the album does not exist.
func (s *Service) GetByTitle(ctx context.Context, title string) (*Album, error) {
	errAlbumWasFound := errors.New("album was found")
	var result *Album
	if err := s.photos.List().ExcludeNonAppCreatedData().PageSize(maxAlbumsPerPage).Pages(ctx, func(response *photoslibrary.ListAlbumsResponse) error {
		if album, found := findByTitle(title, response.Albums); found {
			result = album
			return errAlbumWasFound
		}
		return nil
	}); errors.Is(err, errAlbumWasFound) {
		return result, nil
	}
	return nil, fmt.Errorf("getByTitle %s: %w", title, ErrAlbumNotFound)
}

// List lists all albums in created by this app.
func (s *Service) List(ctx context.Context) ([]Album, error) {
	var result []Album
	albumsListCall := s.photos.List().PageSize(maxAlbumsPerPage).ExcludeNonAppCreatedData()
	err := albumsListCall.Pages(ctx, func(response *photoslibrary.ListAlbumsResponse) error {
		for _, res := range response.Albums {
			result = append(result, toAlbum(res))
		}
		return nil
	})
	return result, err
}

func findByTitle(title string, albums []*photoslibrary.Album) (*Album, bool) {
	for _, a := range albums {
		if a.Title == title {
			album := toAlbum(a)
			return &album, true
		}
	}
	return nil, false
}

func toAlbum(a *photoslibrary.Album) Album {
	return Album{
		ID:                a.Id,
		Title:             a.Title,
		ProductURL:        a.ProductUrl,
		IsWriteable:       a.IsWriteable,
		TotalMediaItems:   a.TotalMediaItems,
		CoverPhotoBaseURL: a.CoverPhotoBaseUrl,
	}
}
