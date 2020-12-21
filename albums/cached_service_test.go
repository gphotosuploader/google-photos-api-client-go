package albums

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestCachedAlbumsService_AddMediaItems(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return success on success", "foo", false},
	}

	r := MockedRepository{
		AddManyItemsFn: func(ctx context.Context, albumId string, mediaItemIds []string) error {
			if "should-fail" == albumId {
				return errors.New("error")
			}
			return nil
		},
	}
	s := CachedAlbumsService{
		repo:  r,
		cache: MockedCache{},
	}
	mediaItems := []string{"foo", "bar", "baz"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.AddMediaItems(context.Background(), tc.input, mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}
func TestCachedAlbumsService_RemoveMediaItems(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return success on success", "foo", false},
	}

	r := MockedRepository{
		RemoveManyItemsFn: func(ctx context.Context, albumId string, mediaItemIds []string) error {
			if "should-fail" == albumId {
				return errors.New("error")
			}
			return nil
		},
	}
	s := CachedAlbumsService{
		repo:  r,
		cache: MockedCache{},
	}
	mediaItems := []string{"foo", "bar", "baz"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.RemoveMediaItems(context.Background(), tc.input, mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}

func TestCachedAlbumsService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return an album on success", "foo", false},
	}

	r := MockedRepository{
		CreateFn: func(ctx context.Context, title string) (*Album, error) {
			if "should-fail" == title {
				return &Album{}, errors.New("error")
			}
			return &Album{Title: title}, nil
		},
	}
	c := MockedCache{
		PutAlbumFn: func(ctx context.Context, album Album) error {
			return nil
		},
	}
	s := CachedAlbumsService{
		repo:  r,
		cache: c,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Create(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.input != got.Title {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func TestCachedAlbumsService_GetById(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return an album on success", "foo", false},
	}

	r := MockedRepository{
		GetFn: func(ctx context.Context, albumId string) (*Album, error) {
			if "should-fail" == albumId {
				return &Album{}, errors.New("error")
			}
			return &Album{ID: albumId}, nil
		},
	}
	c := MockedCache{
		PutAlbumFn: func(ctx context.Context, album Album) error {
			return nil
		},
	}
	s := CachedAlbumsService{
		repo:  r,
		cache: c,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.GetById(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.input != got.ID {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func TestCachedAlbumsService_GetByTitle(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return an album from cache", "album-in-cache", false},
		{"Should return an album on success", "foo", false},
	}

	r := MockedRepository{
		GetByTitleFn: func(ctx context.Context, title string) (*Album, error) {
			if "should-fail" == title {
				return &Album{}, errors.New("error")
			}
			return &Album{Title: title}, nil
		},
	}
	c := MockedCache{
		GetAlbumFn: func(ctx context.Context, title string) (Album, error) {
			if "album-in-cache" == title {
				return Album{Title: title}, nil
			}
			return Album{}, ErrCacheMiss
		},
		PutAlbumFn: func(ctx context.Context, album Album) error {
			return nil
		},
	}
	s := CachedAlbumsService{
		repo:  r,
		cache: c,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.GetByTitle(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.input != got.Title {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func TestCachedAlbumsService_List(t *testing.T) {
	t.Run("Should return all the albums on success", func(t *testing.T) {
		r := MockedRepository{
			ListAllFn: func(ctx context.Context) ([]Album, error) {
				return []Album{
					{Title: "foo"},
					{Title: "bar"},
					{Title: "baz"},
				}, nil
			},
		}
		c := MockedCache{
			InvalidateAllAlbumsFn: func(ctx context.Context) error {
				return nil
			},
			PutManyAlbumsFn: func(ctx context.Context, albums []Album) error {
				return nil
			},
		}

		s := CachedAlbumsService{
			repo:  r,
			cache: c,
		}
		got, err := s.List(context.Background())
		if err != nil {
			t.Fatalf("error was not expected, err: %s", err)
		}
		if 3 != len(got) {
			t.Errorf("want: %d, got: %d", 3, len(got))
		}
	})

	t.Run("Should return error on error", func(t *testing.T) {
		r := MockedRepository{
			ListAllFn: func(ctx context.Context) ([]Album, error) {
				return []Album{}, errors.New("error")
			},
		}
		c := MockedCache{
			InvalidateAllAlbumsFn: func(ctx context.Context) error {
				return nil
			},
			PutManyAlbumsFn: func(ctx context.Context, albums []Album) error {
				return nil
			},
		}

		s := CachedAlbumsService{
			repo:  r,
			cache: c,
		}
		if _, err := s.List(context.Background()); err == nil {
			t.Fatalf("error was expected, but not produced")
		}
	})

	t.Run("Should return error on cache error", func(t *testing.T) {
		r := MockedRepository{
			ListAllFn: func(ctx context.Context) ([]Album, error) {
				return []Album{
					{Title: "foo"},
					{Title: "bar"},
					{Title: "baz"},
				}, nil
			},
		}
		c := MockedCache{
			InvalidateAllAlbumsFn: func(ctx context.Context) error {
				return errors.New("error")
			},
			PutManyAlbumsFn: func(ctx context.Context, albums []Album) error {
				return nil
			},
		}

		s := CachedAlbumsService{
			repo:  r,
			cache: c,
		}
		if _, err := s.List(context.Background()); err == nil {
			t.Fatalf("error was expected, but not produced")
		}
	})
}

func TestNewCachedAlbumsService(t *testing.T) {
	t.Run("WithoutCacheOption", func(t *testing.T) {
		c := MockedCache{
			GetAlbumFn: func(ctx context.Context, title string) (Album, error) {
				return Album{Title: title}, nil
			},
		}
		s := NewCachedAlbumsService(http.DefaultClient, WithCache(c))
		got, err := s.GetByTitle(context.Background(), "foo")
		if err != nil {
			t.Fatalf("error was not expected, err: %s", err)
		}
		if "foo" != got.Title {
			t.Errorf("want: %s, got: %s", "foo", got.Title)
		}
	})

	t.Run("WithoutRepositoryOption", func(t *testing.T) {
		r := MockedRepository{
			AddManyItemsFn: func(ctx context.Context, albumId string, mediaItemIds []string) error {
				return nil
			},
		}
		s := NewCachedAlbumsService(http.DefaultClient, WithRepository(r))
		if err := s.AddMediaItems(context.Background(), "foo", []string{"bar"}); err != nil {
			t.Errorf("want: nil, got: %s", err)
		}
	})
}

func assertExpectedError(errExpected bool, err error, t *testing.T) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
