package albums

import "context"

// MockedRepository mocks the repository.
type MockedRepository struct {
	AddManyItemsFn    func(ctx context.Context, albumId string, mediaItemIds []string) error
	RemoveManyItemsFn func(ctx context.Context, albumId string, mediaItemIds []string) error
	CreateFn          func(ctx context.Context, title string) (*Album, error)
	GetFn             func(ctx context.Context, albumId string) (*Album, error)
	ListAllFn         func(ctx context.Context) ([]Album, error)
	GetByTitleFn      func(ctx context.Context, title string) (*Album, error)
}

// AddManyItems invokes the mock implementation.
func (s MockedRepository) AddManyItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return s.AddManyItemsFn(ctx, albumId, mediaItemIds)
}

// RemoveManyItems invokes the mock implementation.
func (s MockedRepository) RemoveManyItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return s.RemoveManyItemsFn(ctx, albumId, mediaItemIds)
}

// Create invokes the mock implementation.
func (s MockedRepository) Create(ctx context.Context, title string) (*Album, error) {
	return s.CreateFn(ctx, title)
}

// Get invokes the mock implementation.
func (s MockedRepository) Get(ctx context.Context, albumId string) (*Album, error) {
	return s.GetFn(ctx, albumId)
}

// ListAll invokes the mock implementation.
func (s MockedRepository) ListAll(ctx context.Context) ([]Album, error) {
	return s.ListAllFn(ctx)
}

// GetByTitle invokes the mock implementation.
func (s MockedRepository) GetByTitle(ctx context.Context, title string) (*Album, error) {
	return s.GetByTitleFn(ctx, title)
}
