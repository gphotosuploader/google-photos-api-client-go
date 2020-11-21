package albums

import (
	"context"
	"errors"
	"testing"

	duffpl "github.com/duffpl/google-photos-api-client/albums"
)

func TestAlbumRepository_AddManyItems(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return success on success", "foo", false},
	}

	m := MockedAlbumsRepository{
		BatchAddMediaItemsAllFn: func(albumId string, mediaItemIds []string, ctx context.Context) error {
			if "should-fail" == albumId {
				return errors.New("error")
			}
			return nil
		},
	}
	ar := DuffplAlbumRepository{duffplAlbumsClient: m}
	ctx := context.Background()
	mediaItems := []string{"foo", "bar", "baz"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ar.AddManyItems(ctx, tc.input, mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}

func TestAlbumRepository_RemoveManyItems(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return success on success", "foo", false},
	}

	m := MockedAlbumsRepository{
		BatchRemoveMediaItemsAllFn: func(albumId string, mediaItemIds []string, ctx context.Context) error {
			if "should-fail" == albumId {
				return errors.New("error")
			}
			return nil
		},
	}
	ar := DuffplAlbumRepository{duffplAlbumsClient: m}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ar.RemoveManyItems(context.Background(), tc.input, []string{"foo", "bar", "baz"})
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}

func TestAlbumRepository_Create(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return an album on success", "foo", false},
	}

	m := MockedAlbumsRepository{
		CreateFn: func(title string, ctx context.Context) (*duffpl.Album, error) {
			if "should-fail" == title {
				return &duffpl.Album{}, errors.New("error")
			}
			return &duffpl.Album{Title: title}, nil
		},
	}
	ar := DuffplAlbumRepository{duffplAlbumsClient: m}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ar.Create(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.input != got.Title {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func TestAlbumRepository_Get(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return an album on success", "foo", false},
	}

	m := MockedAlbumsRepository{
		GetFn: func(title string, ctx context.Context) (*duffpl.Album, error) {
			if "should-fail" == title {
				return &duffpl.Album{}, errors.New("error")
			}
			return &duffpl.Album{Title: title}, nil
		},
	}
	ar := DuffplAlbumRepository{duffplAlbumsClient: m}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ar.Get(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.input != got.Title {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func TestAlbumRepository_GetByTitle(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
		errExpected   error
	}{
		{"Should return the album on success", "bar", false, nil},
		{"Should return ErrAlbumNotFound if the album does not exist", "non-existent", true, ErrAlbumNotFound},
	}

	m := MockedAlbumsRepository{
		ListAllAsyncFn: func(options *duffpl.AlbumsListOptions, ctx context.Context) (<-chan duffpl.Album, <-chan error) {
			albumsInStorage := []string{"foo", "bar", "baz"}
			albumsC := make(chan duffpl.Album, len(albumsInStorage))
			errorsC := make(chan error)
			go func() {
				defer close(albumsC)
				for _, item := range albumsInStorage {
					albumsC <- duffpl.Album{Title: item}
				}
			}()
			return albumsC, errorsC
		},
	}
	ar := DuffplAlbumRepository{duffplAlbumsClient: m}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ar.GetByTitle(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if tc.errExpected != nil && tc.errExpected != err {
				t.Errorf("err want: %s, err got: %s", tc.errExpected, err)
			}
			if err == nil && tc.input != got.Title {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func TestAlbumRepository_ListAll(t *testing.T) {
	t.Run("Should return all the albums on success", func(t *testing.T) {
		m := MockedAlbumsRepository{
			ListAllFn: func(options *duffpl.AlbumsListOptions, ctx context.Context) ([]duffpl.Album, error) {
				return []duffpl.Album{
					{Title: "foo"},
					{Title: "bar"},
					{Title: "baz"},
				}, nil
			},
		}
		ar := DuffplAlbumRepository{duffplAlbumsClient: m}
		got, err := ar.ListAll(context.Background())
		if err != nil {
			t.Fatalf("error was not expected, err: %s", err)
		}
		if 3 != len(got) {
			t.Errorf("want: %d, got: %d", 3, len(got))
		}
	})

	t.Run("Should return error on error", func(t *testing.T) {
		m := MockedAlbumsRepository{
			ListAllFn: func(options *duffpl.AlbumsListOptions, ctx context.Context) ([]duffpl.Album, error) {
				return []duffpl.Album{}, errors.New("error")
			},
		}
		ar := DuffplAlbumRepository{duffplAlbumsClient: m}
		if _, err := ar.ListAll(context.Background()); err == nil {
			t.Fatalf("error was expected, but not produced")
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
