# Google Photos API client for Go
[![Go Documentation](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos)
[![Build Status](https://cloud.drone.io/api/badges/gphotosuploader/google-photos-api-client-go/status.svg)](https://cloud.drone.io/gphotosuploader/google-photos-api-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/gphotosuploader/google-photos-api-client-go)](https://goreportcard.com/report/github.com/gphotosuploader/google-photos-api-client-go)
[![codebeat badge](https://codebeat.co/badges/c0ab08dd-11b3-406e-bbcc-b9d4a90aedf6)](https://codebeat.co/projects/github-com-gphotosuploader-google-photos-api-client-go-master)
[![codecov](https://codecov.io/gh/gphotosuploader/google-photos-api-client-go/branch/master/graph/badge.svg)](https://codecov.io/gh/gphotosuploader/google-photos-api-client-go)
[![GitHub release](https://img.shields.io/github/release/gphotosuploader/google-photos-api-client-go.svg)](https://github.com/gphotosuploader/google-photos-api-client-go/releases/latest)
[![GitHub](https://img.shields.io/github/license/gphotosuploader/google-photos-api-client-go.svg)](LICENSE)

This package provides a client for using the [Google Photos API](https://godoc.org/google.golang.org/api). Uses the [original Google Photos package](https://github.com/gphotosuploader/googlemirror), that was provided by Google, and [removed](https://code-review.googlesource.com/c/google-api-go-client/+/39951) some time ago. On top of the old package has been extended some other features like uploads, including resumable uploads.


> This project will maintain compatibility with the last two Go major versions published. Currently Go 1.12 and Go 1.13. 
>
## Quick start

Construct a new Google Photos client, then use the various services on the client to access different parts of the Google Photos API. For example:

```go
	import "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"

    // httpClient is an authenticated http.Client. See Authentication below.
	client := gphotos.NewClient(httpClient)
    // get or create a Photos Album with the specified name.
	album, err := GetOrCreateAlbumByName("my-new-album")
	// upload an specified file to an existent Photos Album.
    _, err := client.AddMediaItem(ctx, path, albumID)
```

> NOTE: Using the [context package](https://godoc.org/context), one can easily pass cancellation signals and deadlines to various services of the client for handling a request. In case there is no context available, then `context.Background()` can be used as a starting point.

## Authentication
The gphotos library **does not directly handle authentication**. Instead, when creating a new client, pass an `http.Client` that can handle authentication for you. The easiest and recommended way to do this is using the `golang.org/x/oauth2` library, but you can always use any other library that provides an `http.Client`.

Access to the API requires OAuth client credentials from a Google developers project. This project must have the Library API enabled as described [here](https://developers.google.com/photos/library/guides/get-started).

```go
	import (
        "golang.org/x/oauth2"
        "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
    )
	func main() {
		ctx := context.Background()
		oc := oauth2Config := oauth2.Config{
			ClientID:     "... your application Client ID ...",
			ClientSecret: "... your application Client Secret ...",
			Endpoint:     photos.Endpoint,
			Scopes:       photos.Scopes,
		}
		tc := oc.Client(ctx, "... your user Oauth Token ...")
		client := gphotos.NewClient(tc)
		// look for a Google Photos Album by name
		album, _, err := client.AlbumByName(ctx, "my-album")
	}
```

Note that when using an authenticated Client, all calls made by the client will include the specified OAuth token. Therefore, authenticated clients should almost never be shared between different users. See the oauth2 docs for complete instructions on using that library.

## Limitations
### Rate Limiting
Google Photos imposes a rate limit on all API clients. The quota limit for requests to the Library API is 10,000 requests per project per day. The quota limit for requests to access media bytes (by loading a photo or video from a base URL) is 75,000 requests per project per day.

### Photo storage and quality
All media items uploaded to Google Photos using the API [are stored in full resolution]((https://support.google.com/photos/answer/6220791)) at original quality. **They count toward the user’s storage**.
