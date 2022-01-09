package media_items

import (
	"context"
	"errors"
	"testing"
)

func TestHttpMediaItemsService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		filename      string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return success on success", "foo", false},
	}

	r := MockedRepository{
		CreateManyFn: func(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
			if "should-fail" == mediaItems[0].FileName {
				return []MediaItem{}, errors.New("error")
			}
			return []MediaItem{
				{Filename: mediaItems[0].FileName},
			}, nil
		},
	}

	s := HttpMediaItemsService{repo: r}
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Create(ctx, SimpleMediaItem{FileName: tc.filename})
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.filename != got.Filename {
				t.Errorf("want: %s, got: %s", tc.filename, got.Filename)
			}
		})
	}
}

func TestHttpMediaItemsService_CreateMany(t *testing.T) {
	testCases := []struct {
		name          string
		filenames     []string
		want          int
		isErrExpected bool
	}{
		{"Should return error if API fails", []string{"should-fail", "dummy"}, 0, true},
		{"Should return success on success", []string{"foo", "bar", "baz"}, 3, false},
	}

	r := MockedRepository{
		CreateManyFn: func(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
			ret := make([]MediaItem, len(mediaItems))
			for i, item := range mediaItems {
				if "should-fail" == item.FileName {
					return []MediaItem{}, errors.New("error")
				}
				ret[i].Filename = item.FileName
			}
			return ret, nil
		},
	}

	s := HttpMediaItemsService{repo: r}
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mediaItems := make([]SimpleMediaItem, len(tc.filenames))
			for i, filename := range tc.filenames {
				mediaItems[i].FileName = filename
			}
			got, err := s.CreateMany(ctx, mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.want != len(got) {
				t.Errorf("want: %d, got: %d", tc.want, len(got))
			}
		})
	}
}

func TestHttpMediaItemsService_CreateToAlbum(t *testing.T) {
	testCases := []struct {
		name          string
		albumId       string
		filename      string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", "dummy", true},
		{"Should return success on success", "foo", "bar", false},
	}

	r := MockedRepository{
		CreateManyToAlbumFn: func(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
			if "should-fail" == albumId {
				return []MediaItem{}, errors.New("error")
			}
			return []MediaItem{
				{Filename: mediaItems[0].FileName},
			}, nil
		},
	}

	s := HttpMediaItemsService{repo: r}
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.CreateToAlbum(ctx, tc.albumId, SimpleMediaItem{FileName: tc.filename})
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.filename != got.Filename {
				t.Errorf("want: %s, got: %s", tc.filename, got.Filename)
			}
		})
	}
}

func TestHttpMediaItemsService_CreateManyToAlbum(t *testing.T) {
	testCases := []struct {
		name          string
		albumId       string
		filenames     []string
		want          int
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", []string{"dummy-1", "dummy-2"}, 0, true},
		{"Should return success on success", "foo", []string{"bar", "baz"}, 2, false},
	}

	r := MockedRepository{
		CreateManyToAlbumFn: func(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
			if "should-fail" == albumId {
				return []MediaItem{}, errors.New("error")
			}
			ret := make([]MediaItem, len(mediaItems))
			for i, item := range mediaItems {
				ret[i].Filename = item.FileName
			}
			return ret, nil
		},
	}

	s := HttpMediaItemsService{repo: r}
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mediaItems := make([]SimpleMediaItem, len(tc.filenames))
			for i, filename := range tc.filenames {
				mediaItems[i].FileName = filename
			}
			got, err := s.CreateManyToAlbum(ctx, tc.albumId, mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.want != len(got) {
				t.Errorf("want: %d, got: %d", tc.want, len(got))
			}
		})
	}
}

func TestHttpMediaItemsService_Get(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", "", true},
		{"Should return success on success", "foo", "foo", false},
	}

	r := MockedRepository{
		GetFn: func(ctx context.Context, itemId string) (*MediaItem, error) {
			if "should-fail" == itemId {
				return &MediaItem{}, errors.New("error")
			}
			return &MediaItem{ID: itemId}, nil
		},
	}

	s := HttpMediaItemsService{repo: r}
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Get(ctx, tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.want != got.ID {
				t.Errorf("want: %s, got: %s", tc.want, got.ID)
			}
		})
	}
}

func TestHttpMediaItemsService_ListByAlbum(t *testing.T) {
	testCases := []struct {
		name  string
		input string

		isErrExpected bool
	}{
		{"Should return error if API fails", "should-fail", true},
		{"Should return success on success", "foo", false},
	}

	r := MockedRepository{
		ListByAlbumFn: func(ctx context.Context, albumId string) ([]MediaItem, error) {
			if "should-fail" == albumId {
				return []MediaItem{}, errors.New("error")
			}
			return []MediaItem{
				{ID: "item-1", Filename: "filename-1"},
				{ID: "item-2", Filename: "filename-2"},
			}, nil
		},
	}

	s := HttpMediaItemsService{repo: r}
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.ListByAlbum(ctx, tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && 2 != len(got) {
				t.Errorf("want: %d, got: %d", 2, len(got))
			}
		})
	}
}

func assertExpectedError(isErrExpected bool, err error, t *testing.T) {
	if isErrExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !isErrExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
