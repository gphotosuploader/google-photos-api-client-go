// Package gphotos provides a client for calling the Google Photos API.
//
// Usage:
//    import gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
//
// Construct a new Google Photos client, it needs an authenticated HTTP Client see Authentication section below.
//    client, err := gphotos.NewClient(httpClient)
//    ...
//
// Use WithUploader(), WithAlbumsService(), WithMediaItemsService() to customize it.
// By default:
//     - Uses an Album repository with an in-memory cache.
//     - Uses a basic HTTP uploader, you can find a resumable one using uploader.NewResumableUploader().
//
// It can get Album from the library, returning ErrAlbumNotFound in case it does not exist:
//    title := "my-album"
//    album, err := client.Albums.GetByTitle(ctx, title)
//    if errors.Is(err, ErrAlbumNotFound) {
//       // album does not exist
//    }
//    ...
//
// It can upload a new item to your library:
//    media, err := client.UploadFileToLibrary(ctx, "/my-folder/my-picture.jpg")
//    if err != nil {
//       // handle error
//    }
//    ...
//
// Or upload and adding it to an Album:
//    media, err := client.UploadFileToAlbum(ctx, album.ID, "/my-folder/my-picture.jpg")
//    if err != nil {
//       // handle error
//    }
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
// Limitations
//
// Google Photos API imposes some limitations, please read them all at:
// https://github.com/gphotosuploader/google-photos-api-client-go/
//
package gphotos // import "github.com/gphotosuploader/google-photos-api-client-go/v2"
