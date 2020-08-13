package gphotos_test

import (
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

type mockUploadSessionStore struct{}

func (m *mockUploadSessionStore) Get(f string) []byte {
	return []byte(f)
}

func (m *mockUploadSessionStore) Set(f string, u []byte) {}

func (m *mockUploadSessionStore) Delete(f string) {}

func TestNewClientWithResumableUploads(t *testing.T) {
	c := http.DefaultClient
	store := &mockUploadSessionStore{}

	t.Run("EmptyHTTPClient", func(t *testing.T) {
		_, err := gphotos.NewClientWithResumableUploads(nil, store)
		if err == nil {
			t.Errorf("NewClientWithResumableUploads error was expected here")
		}
	})

	t.Run("WithNilUploadSessionStore", func(t *testing.T) {
		_, err := gphotos.NewClientWithResumableUploads(c, nil)
		if err != uploader.ErrNilStore {
			t.Errorf("NewClientWithResumableUploads - error was expected here: got=%s, want=%s", err, uploader.ErrNilStore)
		}
	})

	t.Run("WithoutOptions", func(t *testing.T) {
		got, err := gphotos.NewClientWithResumableUploads(c, store)
		if err != nil {
			t.Errorf("NewClientWithResumableUploads - error was not expected here: err=%s", err)
		}
		if got.Service == nil {
			t.Errorf("NewClientWithResumableUploads - Photos service was not created")
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		l := log.NewDiscardLogger()
		got, err := gphotos.NewClientWithResumableUploads(c, store, gphotos.WithLogger(l))
		if err != nil {
			t.Errorf("NewClientWithResumableUploads - error was not expected here: err=%s", err)
		}
		if got.Service == nil {
			t.Errorf("NewClientWithResumableUploads - Photos service was not created")
		}
	})
}

// DEPRECATED
func TestNewClient(t *testing.T) {
	want := http.DefaultClient

	t.Run("WithoutEmptyHTTPClient", func(t *testing.T) {
		_, err := gphotos.NewClient(nil)
		if err == nil {
			t.Errorf("error was expected here")
		}
	})

	t.Run("WithoutToken", func(t *testing.T) {
		got, err := gphotos.NewClient(want)
		if err != nil {
			t.Errorf("error was not expected here: err=%s", err)
		}

		if got.Service == nil {
			t.Errorf("Photos service was not created")
		}
	})

	t.Run("WithToken", func(t *testing.T) {
		tk := testOauthToken()
		got, err := gphotos.NewClient(want, &tk)
		if err != nil {
			t.Errorf("error was not expected here: err=%s", err)
		}

		if got.Service == nil {
			t.Errorf("Photos service was not created")
		}

		if *(got.Token()) != tk {
			t.Errorf("Token is different from expected")
		}
	})
}

// DEPRECATED
func TestClient_Token(t *testing.T) {
	c := http.DefaultClient

	t.Run("EmptyToken", func(t *testing.T) {
		got, err := gphotos.NewClient(c, nil)
		if err != nil {
			t.Errorf("error was not expected here: err=%s", err)
		}

		if got.Token() != nil {
			t.Errorf("Token should be nil: got:%v", got.Token())
		}
	})

	t.Run("ValidToken", func(t *testing.T) {
		tk := testOauthToken()
		got, err := gphotos.NewClient(c, &tk)
		if err != nil {
			t.Errorf("error was not expected here: err=%s", err)
		}

		if *(got.Token()) != tk {
			t.Errorf("Token is different from expected")
		}
	})
}

// DEPRECATED
func testOauthToken() oauth2.Token {
	return oauth2.Token{
		AccessToken:  "access-token",
		TokenType:    "token-type",
		RefreshToken: "refresh-token",
		Expiry:       time.Time{},
	}
}
