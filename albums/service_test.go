package albums_test

import (
	"context"
	"errors"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
	"net/http"
	"testing"
)

func TestAlbumsService_AddMediaItems(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return success on success", "foo", false},
	}

	mockedRepository := albums.MockedRepository{
		AddManyItemsFn: func(ctx context.Context, albumID string, mediaItemIds []string) error {
			if "should-fail" == albumID {
				return errors.New("error")
			}
			return nil
		},
	}

	s := albums.NewService(
		http.DefaultClient,
		albums.WithRepository(mockedRepository),
	)

	mediaItems := []string{"foo", "bar", "baz"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.AddMediaItems(context.Background(), tc.input, mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}
func TestAlbumsService_RemoveMediaItems(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return success on success", "foo", false},
	}

	mockedRepository := albums.MockedRepository{
		RemoveManyItemsFn: func(ctx context.Context, albumID string, mediaItemIds []string) error {
			if "should-fail" == albumID {
				return errors.New("error")
			}
			return nil
		},
	}

	s := albums.NewService(
		http.DefaultClient,
		albums.WithRepository(mockedRepository),
	)

	mediaItems := []string{"foo", "bar", "baz"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.RemoveMediaItems(context.Background(), tc.input, mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}

func TestAlbumsService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return an album on success", "foo", false},
	}

	mockedRepository := albums.MockedRepository{
		CreateFn: func(ctx context.Context, title string) (*albums.Album, error) {
			if "should-fail" == title {
				return &albums.NullAlbum, errors.New("error")
			}
			return &albums.Album{Title: title}, nil
		},
	}

	s := albums.NewService(
		http.DefaultClient,
		albums.WithRepository(mockedRepository),
	)

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

func TestAlbumsService_GetById(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return an album on success", "foo", false},
	}

	mockedRepository := albums.MockedRepository{
		GetFn: func(ctx context.Context, albumId string) (*albums.Album, error) {
			if "should-fail" == albumId {
				return &albums.NullAlbum, errors.New("error")
			}
			return &albums.Album{ID: albumId}, nil
		},
	}

	s := albums.NewService(
		http.DefaultClient,
		albums.WithRepository(mockedRepository),
	)

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

func TestAlbumsService_GetByTitle(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return an album on success", "foo", false},
	}

	mockedRepository := albums.MockedRepository{
		GetByTitleFn: func(ctx context.Context, title string) (*albums.Album, error) {
			if "should-fail" == title {
				return &albums.NullAlbum, errors.New("error")
			}
			return &albums.Album{Title: title}, nil
		},
	}

	s := albums.NewService(
		http.DefaultClient,
		albums.WithRepository(mockedRepository),
	)

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

func TestAlbumsService_List(t *testing.T) {
	t.Run("Should return all the albums on success", func(t *testing.T) {
		mockedRepository := albums.MockedRepository{
			ListAllFn: func(ctx context.Context) ([]albums.Album, error) {
				return []albums.Album{
					{Title: "foo"},
					{Title: "bar"},
					{Title: "baz"},
				}, nil
			},
		}

		s := albums.NewService(
			http.DefaultClient,
			albums.WithRepository(mockedRepository),
		)

		got, err := s.List(context.Background())
		if err != nil {
			t.Fatalf("error was not expected, err: %s", err)
		}
		if 3 != len(got) {
			t.Errorf("want: %d, got: %d", 3, len(got))
		}
	})

	t.Run("Should return error on error", func(t *testing.T) {
		mockedRepository := albums.MockedRepository{
			ListAllFn: func(ctx context.Context) ([]albums.Album, error) {
				return []albums.Album{}, errors.New("error")
			},
		}

		s := albums.NewService(
			http.DefaultClient,
			albums.WithRepository(mockedRepository),
		)

		if _, err := s.List(context.Background()); err == nil {
			t.Fatalf("error was expected, but not produced")
		}
	})
}
