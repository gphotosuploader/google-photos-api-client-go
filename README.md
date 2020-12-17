# Google Photos API client for Go
[![Go Documentation](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/gphotosuploader/google-photos-api-client-go/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/gphotosuploader/google-photos-api-client-go)](https://goreportcard.com/report/github.com/gphotosuploader/google-photos-api-client-go)
[![codebeat badge](https://codebeat.co/badges/c0ab08dd-11b3-406e-bbcc-b9d4a90aedf6)](https://codebeat.co/projects/github-com-gphotosuploader-google-photos-api-client-go-master)
[![codecov](https://codecov.io/gh/gphotosuploader/google-photos-api-client-go/branch/master/graph/badge.svg)](https://codecov.io/gh/gphotosuploader/google-photos-api-client-go)
[![GitHub release](https://img.shields.io/github/release/gphotosuploader/google-photos-api-client-go.svg)](https://github.com/gphotosuploader/google-photos-api-client-go/releases/latest)
[![GitHub](https://img.shields.io/github/license/gphotosuploader/google-photos-api-client-go.svg)](LICENSE)

This package provides a client for using the [Google Photos API](https://godoc.org/google.golang.org/api) in Go. Uses the [original Google Photos package](https://github.com/gphotosuploader/googlemirror), that was provided by Google, and [removed](https://code-review.googlesource.com/c/google-api-go-client/+/39951) some time ago. 
 
> This project will maintain compatibility with the last two Go major versions [published](https://golang.org/doc/devel/release.html). 

## Quick start

```go
package main

import gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"

func main() error {
	ctx := context.Background()

	// httpClient is an authenticated http.Client. See Authentication below. 
	// You can customize the client using Options, see below.
	client, err := gphotos.NewClient(httpClient)
	if err != nil {
		// handle error
	}

	// create a Photos Album with the specified name.
	title := "my-album"
	album, err := client.CreateAlbum(ctx, title)
	if err != nil {
		// handle error
	}

	// upload an specified file to the previous album.
	item := ghotos.FileUploadItem("/my-folder/my-picture.jpg")
	_, err = client.AddMediaToAlbum(ctx, item, album)
	if err != nil {
		// handle error
	}   
	...
}
```

## Features

* **Two uploaders**: We have implemented the `/v1/uploads` endpoint to be able to upload media items to your library. Two implementations are available: a `BasicUploader` and a `ResumableUploader` which could be used for large files, like videos. See [Options](#options) 
* **Uploads**: `AddMediaToLibrary` and `AddMediaToAlbum` allows you to upload media to the the libray.
* **Album management**: `ListAlbums`, `FindAlbum` and `CreateAlbum` allows you to list, find an album by title and create a new album in the library. 
* `ListAlbumsWithCallback` allows you to call a function while browse the albums in the library.
* **Automatic cache**, following [Google Photos best practices](https://developers.google.com/photos/library/guides/best-practices#caching), a cache is used to increase performance and reduce the number of calls to the Google Photos API (see [Rate Limiting](#rate-limiting).
* **Automatic retries**: It implements retries following [Google Photos error handling documentation](https://developers.google.com/photos/library/guides/best-practices#error-handling). By default it uses exponential backoff with a maximum of 5 retries.

## Options 

There are several options to customize the Google Photos client:

* `WithPhotoService`: Allows you to use a different Google Photos service. It should implement [photoservice.Service](internal/photoservice/types.go) interface. By default it's using [this one](https://github.com/gphotosuploader/googlemirror).  

* `WithUploader`: Allows you to use your own uploader. It should implement [uploader.Uploader](internal/uploader/uploader.go) interface. By default it's using [BasicUploader](internal/uploader/basic.go). See `WithSessionStorer` to use [ResumableUploader](internal/uploader/resumable.go).

* `WithSessionStorer`: Enables a [ResumableUploader](internal/uploader/resumable.go) using the provided storage to keep upload sessions (see [uploader.SessionStorer](internal/uploader/resumable.go) interface).

* `WithCacher`: Allows you to use your own cache. It should implement [cache.Cache](albums/cache/cache.go) interface. By default it's using `gadelkareem/cachita`.

* `WithLogger`: Allows you to use a logger to print messages from this library. It should implement [log.Logger](internal/log/logger.go). By default, there is no logging.

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
	}
```

Note that when using an authenticated Client, all calls made by the client will include the specified OAuth token. Therefore, authenticated clients should almost never be shared between different users. See the oauth2 docs for complete instructions on using that library.

## Limitations
### Rate Limiting
Google Photos imposes a [rate limit on all API clients](https://developers.google.com/photos/library/guides/api-limits-quotas). The quota limit for requests to the Library API is 10,000 requests per project per day. The quota limit for requests to access media bytes (by loading a photo or video from a base URL) is 75,000 requests per project per day.

### Photo storage and quality
All media items uploaded to Google Photos using the API [are stored in full resolution](https://developers.google.com/photos/library/guides/api-limits-quotas) at original quality. **They count toward the user’s storage**.

### Uploads are only allowed to albums created by this application
Note that you can only add media items that have been uploaded by this application to albums that this application has created, see [here](https://developers.google.com/photos/library/guides/manage-albums#adding-items-to-album) why.
