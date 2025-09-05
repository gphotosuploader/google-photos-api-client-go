# Google Photos API client for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/gphotosuploader/google-photos-api-client-go/v3.svg)](https://pkg.go.dev/github.com/gphotosuploader/google-photos-api-client-go/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/gphotosuploader/google-photos-api-client-go)](https://goreportcard.com/report/github.com/gphotosuploader/google-photos-api-client-go)
[![codecov](https://codecov.io/gh/gphotosuploader/google-photos-api-client-go/branch/main/graph/badge.svg)](https://codecov.io/gh/gphotosuploader/google-photos-api-client-go)
[![GitHub release](https://img.shields.io/github/release/gphotosuploader/google-photos-api-client-go.svg)](https://github.com/gphotosuploader/google-photos-api-client-go/releases/latest)
[![GitHub](https://img.shields.io/github/license/gphotosuploader/google-photos-api-client-go.svg)](LICENSE)

[iDocumentation]: https://pkg.go.dev/github.com/gphotosuploader/google-photos-api-client-go/v2

This package provides a client for using the [Google Photos API](https://developers.google.com/photos) in Go. Uses the original `photoslibrary` package, that was [provided by Google](https://code-review.googlesource.com/c/google-api-go-client/+/39951) and, now it's archived [here](https://github.com/gphotosuploader/googlemirror).

The package offers access to these Google Photos services:
- `Albums` is a service to manage albums.
- `MediaItems` is a service to manage media items (Photos and Videos).
- `Uploader` is a service to upload items.

> This project will maintain compatibility with the last three major [published](https://golang.org/doc/devel/release.html) versions of Go.

## Installation

`google-photos-api-client-go` is compatible with modern Go releases in module mode. You can install it using:

```bash
$ go get github.com/gphotosuploader/google-photos-api-client-go/v3
```

## Usage

```go
import gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
```

Construct a new Google Photos client, then use the various services on the client to access different parts of the Google Photos API. 

For example:

```go
// httpClient has been previously authenticated using oAuth authenticated
client := gphotos.NewClient(httpClient)

// list all albums for the authenticated user
albums, err := client.Albums.List(context.Background())
```

The services of a client divide the API into logical chunks and correspond to the structure of the Google Photos API documentation at https://developers.google.com/photos/library/reference/rest.

**NOTE**: Using the [context](https://godoc.org/context) package, one can easily pass cancellation signals and deadlines to various services of the client for handling a request. In case there is no context available, then `context.Background()` can be used as a starting point.

## Authentication
The gphotos library **does not directly handle authentication**. Instead, when creating a new client, pass a `net/http.Client` that can handle authentication for you. The easiest and recommended way to do this is using the [oauth2](https://github.com/golang/oauth2) library, but you can always use any other library that provides a `net/http.Client`.

Access to the API requires OAuth 2.0 client credentials from a Google developers project. This project must have the Library API enabled as described [here](https://developers.google.com/photos/library/guides/get-started).

```
import (
    "golang.org/x/oauth2"

    gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
)

func main() {
    ctx := context.Background()
    oc := oauth2Config := oauth2.Config{
        ClientID:     "... your application Client ID ...",
        ClientSecret: "... your application Client Secret ...",
        // ...
    }
    tc := oc.Client(ctx, "... your user Oauth Token ...")
    
    client := gphotos.NewClient(tc)
    
    // list all albums for the authenticated user
    albums, err := client.Albums.List(ctx)
}
```

Note that when using an authenticated Client, all calls made by the client will include the specified OAuth token. Therefore, authenticated clients should almost never be shared between different users.

See the [oauth2 docs](https://godoc.org/golang.org/x/oauth2) for complete instructions on using that library.

## Features

### Retries following best practices

- Follows [Google Photos error handling best practices](https://developers.google.com/photos/library/guides/best-practices#error-handling), using an exponential backoff retrier with a maximum of 3 retries.
- Returns a specific error type, `ErrDailyQuotaExceeded`, if the 'All requests' per day quota has been exceeded. See [Rate Limiting](#rate-limiting).

### Albums service

- Offers an independent `albums.Service` implementing the [Google Photos Albums API](https://developers.google.com/photos/library/reference/rest#rest-resource:-v1.albums).
- The client accepts a customized albums service using `client.Albums`.
- Consider implementing a [caching strategy](https://developers.google.com/photos/library/guides/best-practices#caching) to avoid [Rate Limiting](#rate-limiting).

### Media Items service

- Offers an independent `albums.Service` implementing the [Google Photos MediaItems API](https://developers.google.com/photos/library/reference/rest#rest-resource:-v1.mediaitems).
- The client accepts a customized media items service using `client.MediaItems`.

### Uploader

- Offers **two upload clients** implementing the [Google Photos Uploads API](https://developers.google.com/photos/library/guides/upload-media).
    - `uploader.SimpleUploader` is a simple HTTP uploader.
    - `uploader.ResumableUploader` is an uploader implementing resumable uploads. It could be used for large files, like videos. See [documentation](https://developers.google.com/photos/library/guides/resumable-uploads).
- The client accepts a customized media items service using `client.Uploader`.

## Limitations
Only images and videos can be uploaded. If you attempt to upload non-videos or images or formats that Google Photos doesn't understand, Google Photos will give an error when creating media item.

### Photo storage and quality
All media items uploaded to Google Photos using the API [are stored in full resolution](https://support.google.com/photos/answer/6220791) at original quality. **They count toward the user’s storage**. The API does not offer a way to upload in "high quality" mode.

### Duplicates
If you upload the same image twice, with the same binary content, then Google Photos will deduplicate it. However, it will retain the filename from the first upload which may be confusing. In practice, this shouldn't cause too many problems.

### Albums
Note that you can only add media items that have been uploaded by this application to albums that this application has created, see [here](https://developers.google.com/photos/library/guides/manage-albums#adding-items-to-album) why.

### Rate Limiting
Google Photos imposes a rate limit on all API clients. **The quota limit for requests to the Library API is 10,000 requests per project per day**. The quota limit for requests to access media bytes (by loading a photo or video from a base URL) is 75,000 requests per project per day.

## Used by

* [gphotos-uploader-cli](https://github.com/gphotosuploader/gphotos-uploader-cli): A command line to sync your pictures and videos with Google Photos. Supporting linux/macOS.
* [Send To Google Photos app](https://github.com/arran4/send-to-google-photos): A simple "Send To" extension for sending images to Google Photos
