package shared_albums

import (
	"context"
	"net/http"

	"github.com/duffpl/google-photos-api-client/shared_albums"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
)

var (
	// Excludes non app created albums. Google Photos doesn't allow manage non created shared_albums through the API.
	// https://developers.google.com/photos/library/guides/manage-albums#adding-items-to-album
	excludeNonAppCreatedData = &ListOptions{ExcludeNonAppCreatedData: true}
)

// SharedAlbumsService represents a Google Photos client for shared albums management.
type SharedAlbumsService interface {
	Get(ctx context.Context, shareToken string) (*albums.Album, error)
	Join(ctx context.Context, shareToken string) (*albums.Album, error)
	Leave(ctx context.Context, shareToken string) error
	List(ctx context.Context) ([]albums.Album, error)
}

type ListOptions = shared_albums.ListOptions

// repository represents the Google Photos API client.
type repository interface {
	Get(shareToken string, ctx context.Context) (*albums.Album, error)
	Join(shareToken string, ctx context.Context) (*albums.Album, error)
	Leave(shareToken string, ctx context.Context) error
	ListAll(options *ListOptions, ctx context.Context) ([]albums.Album, error)
}

func defaultRepo(authenticatedClient *http.Client) shared_albums.HttpSharedAlbumsService {
	return shared_albums.NewHttpSharedAlbumsService(authenticatedClient)
}

// CachedAlbumsService implements a Google Photos client with cached results.
type HttpSharedAlbumsService struct {
	repo repository
}

// Get fetches album based on specified shareToken.
func (s HttpSharedAlbumsService) Get(ctx context.Context, shareToken string) (*albums.Album, error) {
	return s.repo.Get(shareToken, ctx)
}

// Join joins a shared album on behalf of the Google Photos user.
func (s HttpSharedAlbumsService) Join(ctx context.Context, shareToken string) (*albums.Album, error) {
	return s.repo.Join(shareToken, ctx)
}

// Leave leaves a previously-joined shared album on behalf of the Google Photos user. The user must not own this album.
func (s HttpSharedAlbumsService) Leave(ctx context.Context, shareToken string) error {
	return s.repo.Leave(shareToken, ctx)
}

// List lists all shared albums available in the Sharing tab of the user's Google Photos app.
func (s HttpSharedAlbumsService) List(ctx context.Context) ([]albums.Album, error) {
	return s.repo.ListAll(excludeNonAppCreatedData, ctx)
}

// NewSharedAlbumsService returns a client of HttpSharedAlbumsService.
func NewSharedAlbumsService(authenticatedClient *http.Client, options ...Option) HttpSharedAlbumsService {
	var repo repository = defaultRepo(authenticatedClient)

	for _, o := range options {
		switch o.Name() {
		case optkeyRepo:
			repo = o.Value().(repository)
		}
	}

	return HttpSharedAlbumsService{
		repo:  repo,
	}
}

const (
	optkeyRepo  = "repository"
)

// Option represents a configurable parameter.
type Option interface {
	Name() string
	Value() interface{}
}

type option struct {
	name  string
	value interface{}
}

func (o option) Name() string       { return o.name }
func (o option) Value() interface{} { return o.value }

// WithRepository configures the Google Photos repository.
func WithRepository(s repository) Option {
	return &option{
		name:  optkeyRepo,
		value: s,
	}
}
