package gphotos_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
)

func TestNewClientWithOptions(t *testing.T) {
	c := http.DefaultClient

	t.Run("EmptyHTTPClient", func(t *testing.T) {
		_, err := gphotos.NewClientWithOptions(nil)
		if err == nil {
			t.Errorf("NewClientWithOptions error was expected here")
		}
	})

	t.Run("WithoutOptions", func(t *testing.T) {
		got, err := gphotos.NewClientWithOptions(c)
		if err != nil {
			t.Errorf("NewClientWithOptions - error was not expected here: err=%s", err)
		}
		if got.Service == nil {
			t.Errorf("NewClientWithOptions - Photos service was not created")
		}
		if got.Client != c {
			t.Errorf("NewClientWithOptions - HTTP Client is not the expected")
		}
	})

	t.Run("WithoutOptionResume", func(t *testing.T) {
		s := gphotos.Store()
		got, err := gphotos.NewClientWithOptions(c, gphotos.OptionResumeUploads(s))
		if err != nil {
			t.Errorf("NewClientWithOptions - error was not expected here: err=%s", err)
		}
		if got.Service == nil {
			t.Errorf("NewClientWithOptions - Photos service was not created")
		}
		if got.Client != c {
			t.Errorf("NewClientWithOptions - HTTP Client is not the expected")
		}
		if !got.CanResumeUploads() {
			t.Errorf("NewClientWithOptions - got=%t, want=true", got.CanResumeUploads())
		}
	})

	t.Run("WithoutOptionLog", func(t *testing.T) {
		l := log.New(ioutil.Discard, "", 0)
		got, err := gphotos.NewClientWithOptions(c, gphotos.OptionLog(l))
		if err != nil {
			t.Errorf("NewClientWithOptions - error was not expected here: err=%s", err)
		}
		if got.Service == nil {
			t.Errorf("NewClientWithOptions - Photos service was not created")
		}
		if got.Client != c {
			t.Errorf("NewClientWithOptions - HTTP Client is not the expected")
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

		if got.Client != want {
			t.Errorf("HTTP Client is different")
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

		if got.Client != want {
			t.Errorf("HTTP Client is different from expected")
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
