package mock

import (
	"context"

	"github.com/duffpl/google-photos-api-client/albums"
	"github.com/duffpl/google-photos-api-client/shared_albums"
)

// SharedAlbumsRepository mocks the a shared_albums repository.
type SharedAlbumsRepository struct {
	GetFn      func(shareToken string, ctx context.Context) (*albums.Album, error)
	GetInvoked bool

	ListAllFn      func(options *shared_albums.ListOptions, ctx context.Context) ([]albums.Album, error)
	ListAllInvoked bool

	JoinFn      func(shareToken string, ctx context.Context) (*albums.Album, error)
	JoinInvoked bool

	LeaveFn      func(shareToken string, ctx context.Context) error
	LeaveInvoked bool
}

// Get invokes the mock implementation and marks the function as invoked.
func (s SharedAlbumsRepository) Get(shareToken string, ctx context.Context) (*albums.Album, error) {
	s.GetInvoked = true
	return s.GetFn(shareToken, ctx)
}

// ListAll invokes the mock implementation and marks the function as invoked.
func (s SharedAlbumsRepository) ListAll(options *shared_albums.ListOptions, ctx context.Context) ([]albums.Album, error) {
	s.ListAllInvoked = true
	return s.ListAllFn(options, ctx)
}

// Join invokes the mock implementation and marks the function as invoked.
func (s SharedAlbumsRepository) Join(shareToken string, ctx context.Context) (*albums.Album, error) {
	s.JoinInvoked = true
	return s.JoinFn(shareToken, ctx)
}

// Leave invokes the mock implementation and marks the function as invoked.
func (s SharedAlbumsRepository) Leave(shareToken string, ctx context.Context) error {
	s.LeaveInvoked = true
	return s.LeaveFn(shareToken, ctx)
}
