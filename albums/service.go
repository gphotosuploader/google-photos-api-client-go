package albums

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// Repository represents an album repository.
type Repository interface {
	AddManyItems(ctx context.Context, albumId string, mediaItemIds []string) error
	RemoveManyItems(ctx context.Context, albumId string, mediaItemIds []string) error
	Create(ctx context.Context, title string) (*Album, error)
	Get(ctx context.Context, albumId string) (*Album, error)
	ListAll(ctx context.Context) ([]Album, error)
	GetByTitle(ctx context.Context, title string) (*Album, error)
}

// Service implements an albums Google Photos client.
type Service struct {
	repo Repository
}

var (
	// NullAlbum is a zero value Album.
	NullAlbum = Album{}

	// ErrAlbumNotFound is the error returned when an album is not found.
	ErrAlbumNotFound = errors.New("album not found")
)

// AddMediaItems adds multiple media item(s) to the specified album.
func (s Service) AddMediaItems(ctx context.Context, albumID string, mediaItemIDs []string) error {
	return s.repo.AddManyItems(ctx, albumID, mediaItemIDs)
}

// RemoveMediaItems removes multiple media item(s) from the specified album.
func (s Service) RemoveMediaItems(ctx context.Context, albumID string, mediaItemIDs []string) error {
	return s.repo.RemoveManyItems(ctx, albumID, mediaItemIDs)
}

// Create adds a new album to the repo.
func (s Service) Create(ctx context.Context, title string) (*Album, error) {
	return s.repo.Create(ctx, title)
}

// GetById fetches an album from the repository by id.
// It returns ErrAlbumNotFound if the album does not exist.
func (s Service) GetById(ctx context.Context, albumID string) (*Album, error) {
	album, err := s.repo.Get(ctx, albumID)
	if err != nil {
		return &NullAlbum, fmt.Errorf("%s: %w", albumID, ErrAlbumNotFound)
	}
	return album, nil
}

// GetByTitle fetches an album from the repository by title.
// Ir returns ErrAlbumNotFound if the album does not exist.
func (s Service) GetByTitle(ctx context.Context, title string) (*Album, error) {
	album, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		return &NullAlbum, fmt.Errorf("%s: %w", title, ErrAlbumNotFound)
	}
	return album, nil
}

// List fetches all the albums from the repository.
func (s Service) List(ctx context.Context) ([]Album, error) {
	return s.repo.ListAll(ctx)
}

// NewService returns an albums Google Photos client.
// The authenticatedClient should have all oAuth credentials in place.
func NewService(authenticatedClient *http.Client, options ...Option) *Service {
	s := &Service{
		repo: defaultRepo(authenticatedClient),
	}

	for _, o := range options {
		o(s)
	}

	return s
}

type Option func(service *Service)

// WithRepository configures the Google Photos repository.
func WithRepository(repository Repository) Option {
	return func(s *Service) {
		s.repo = repository
	}
}

func defaultRepo(authenticatedClient *http.Client) Repository {
	r, _ := NewPhotosLibraryClient(authenticatedClient)
	return r
}
