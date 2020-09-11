// Package gphotos provides a client for calling the Google Photos API.
//
// Wraps the photoslibray package provided originally by Google, and now maintained in: https://github.com/gphotosuploader/googlemirror.
//
// Usage:
//    import gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
//
// Construct a new Google Photos client, it needs an authenticated HTTP Client see Authentication section below.
//    client, err := NewClient(httpClient)
//    ...
//
// It can get Album from the library, returning ErrAlbumNotFound in case it does not exist:
//    title := "my-album"
//    album, err := client.FindAlbum(ctx, title)
//    if errors.Is(err, ErrAlbumNotFound) {
//       // album does not exist
//    }
//    ...
//
// It can upload a new item to your library:
//    item := FileUploadItem("/my-folder/my-picture.jpg")
//    media, err = client.AddMediaToLibrary(ctx, item)
//    ...
//
// Or upload and adding it to an Album:
//    album, err := client.FindAlbum(ctx, title)
//    if err != nil {
//       // handle error
//    }
//    item := FileUploadItem("/my-folder/my-picture.jpg")
//    media, err = client.AddMediaToLibrary(ctx, item, album)
//    ...
//
// Authentication
//
// The gphotos library does not directly handle authentication. Instead, when
// creating a new client, pass an http.Client that can handle authentication for
// you. The easiest and recommended way to do this is using the golang.org/x/oauth2
// library, but you can always use any other library that provides an http.Client.
// Access to the API requires OAuth client credentials from a Google developers
// project. This project must have the Library API enabled as described in
// https://developers.google.com/photos/library/guides/get-started.
//
//	import "golang.org/x/oauth2"
//	func main() {
//		ctx := context.Background()
//		oc := oauth2Config := oauth2.Config{
//			ClientID:     "... your application Client ID ...",
//			ClientSecret: "... your application Client Secret ...",
//			Endpoint:     photos.Endpoint,
//			Scopes:       photos.Scopes,
//		}
//		tc := oc.Client(ctx, "... your user Oauth Token ...")
//		client, err := gphotos.NewClient(tc)
//      ...
//	}
//
// Note that when using an authenticated Client, all calls made by the client will
// include the specified OAuth token. Therefore, authenticated clients should
// almost never be shared between different users.
// See the oauth2 docs for complete instructions on using that library.
//
// Rate Limiting
//
// Google Photos imposes a rate limit on all API clients. The quota limit for
// requests to the Library API is 10,000 requests per project per day. The quota
// limit for requests to access media bytes (by loading a photo or video from a base
// URL) is 75,000 requests per project per day.
//
// Cached
//
// To improve performance and reduce number of requests to Google Photos, this module
// uses a cache to temporary store some information.
//
// Photo storage and quality
//
// All media items uploaded to Google Photos using the API are stored in full
// resolution at original quality (https://support.google.com/photos/answer/6220791).
// They count toward the user’s storage.
//
package gphotos // import "github.com/gphotosuploader/google-photos-api-client-go/v2"
