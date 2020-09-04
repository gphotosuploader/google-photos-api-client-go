package gphotos_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"google.golang.org/api/googleapi"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/mock"
)

func TestClient_AddMediaToLibrary(t *testing.T) {
	httpClient := http.DefaultClient
	s := &mock.PhotoService{
		CreateMediaItemsFn: func(ctx context.Context, request *photoslibrary.BatchCreateMediaItemsRequest) (*photoslibrary.BatchCreateMediaItemsResponse, error) {
			if request.NewMediaItems[0].Description == "should-fail-on-media-creation" {
				return nil, errors.New("error")
			}

			res := &photoslibrary.BatchCreateMediaItemsResponse{
				ServerResponse: googleapi.ServerResponse{HTTPStatusCode: http.StatusOK},
			}
			res.NewMediaItemResults = append(res.NewMediaItemResults, &photoslibrary.NewMediaItemResult{
				MediaItem: &photoslibrary.MediaItem{
					Description:    request.NewMediaItems[0].Description,
					ServerResponse: googleapi.ServerResponse{HTTPStatusCode: http.StatusOK},
				},
				Status: &photoslibrary.Status{Code: 0},
			})

			return res, nil
		},
	}

	u := &mock.Uploader{
		UploadFn: func(ctx context.Context, item uploader.UploadItem) (token uploader.UploadToken, err error) {
			if item.Name() == "should-fail-on-upload" {
				return "", errors.New("error")
			}
			return "my-token", nil
		},
	}

	c, err := gphotos.NewClient(httpClient, gphotos.WithPhotoService(s), gphotos.WithUploader(u))
	if err != nil {
		t.Errorf("error was not expected at this point")
	}

	t.Run("ReturnsMediaOnSuccess", func(t *testing.T) {
		want := "dummy"
		item := mock.FileUploadItem{Path: want}
		media, err := c.AddMediaToLibrary(context.Background(), item)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
		if media.Description != want {
			t.Errorf("want: %s, got: %s", want, media.Description)
		}
	})

	t.Run("RaiseErrorWhenCreateMediaFails", func(t *testing.T) {
		item := mock.FileUploadItem{Path: "should-fail-on-media-creation"}
		media, err := c.AddMediaToLibrary(context.Background(), item)
		if err == nil {
			t.Errorf("should fail: %v", media)
		}
	})

	t.Run("RaiseErrorWhenUploadFails", func(t *testing.T) {
		item := mock.FileUploadItem{Path: "should-fail-on-upload"}
		media, err := c.AddMediaToLibrary(context.Background(), item)
		if err == nil {
			t.Errorf("should fail, due to upload: %v", media)
		}
	})

}
