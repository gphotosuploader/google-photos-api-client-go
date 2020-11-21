package gphotos

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := NewClient(http.DefaultClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})
}

/*func TestClient_UploadFileToLibrary(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if upload fails", "upload-should-fail", true},
		{"Should return success on success", "foo", false},
	}

	c := Client{
		MediaItems: media_items.MockedMediaItemsService{},
		Uploader:   MockedUploader{
			UploadFileFn: func(ctx context.Context, filepath string) (string, error) {
				if "upload-should-fail" == filepath {
					return "", errors.New("error")
				}
				return "uploadToken", nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := c.UploadFileToLibrary(context.Background(), tc.input)
		})
	}
}*/
