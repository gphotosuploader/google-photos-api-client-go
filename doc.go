// Package gphotos provides a client for calling the [Google Photos API].
//
// # Get started
//
// To start using the [Google Photos API] with this Golang client library, you will need to set up the client library
// in your development environment. Before you do that, configure your project by enabling the API via the Google API
// Console and setting up an OAuth 2.0 client ID. Please follow [this guide] to set it up properly.
//
// Your application interacts with Google Photos on behalf of a Google Photos user. For instance, when you create
// albums in a user's Google Photos library or upload media items to a user's Google Photos account, the user
// authorizes these API requests via the OAuth 2.0 protocol.
//
// The OAuth 2.0 client ID allows your application users to sign in, authenticate, and thereby use the Library API.
//
// # Usage
//
//	import gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
//
// Construct a new Google Photos client using the OAuth 2.0 authenticated HTTP Client, see the Authentication section
// below.
//
//	config := gphotos.Config{Client: httpClient)
//	client, err := gphotos.New(config)
//	...
//
// By default:
//   - Uses a basic HTTP uploader, you can find a resumable one using resumable.NewResumableUploader().
//   - Implements an HTTP retry policy based on the Google Photos best practices.
//
// You can customize the client using the [Config] struct.
//
// It can get an album from the library:
//
//	title := "my-album"
//	album, err := client.Albums.GetByTitle(ctx, title)
//	if err != nil {
//	   // handle error
//	}
//	...
//
// It can upload a new item to your library:
//
//	media, err := client.UploadFileToLibrary(ctx, "/my-folder/my-picture.jpg")
//	if err != nil {
//	   // handle error
//	}
//	...
//
// Or upload and add it to an album:
//
//	media, err := client.UploadFileToAlbum(ctx, album.ID, "/my-folder/my-picture.jpg")
//	if err != nil {
//	   // handle error
//	}
//	...
//
// # Authentication
//
// The gphotos package does not directly handle authentication. Instead, when creating a new client, pass a
// [http.Client] that can handle authentication for you.
//
// The easiest and recommended way to do this is using the [golang.org/x/oauth2] package, but you can always use any
// other package that provides a [net/http.Client].
//
// Access to the API requires an OAuth 2.0 client.
//
//	  import (
//	    "golang.org/x/oauth2"
//
//	    gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
//	   )
//
//	   func main() {
//	     ctx := context.Background()
//	     oc := oauth2Config := oauth2.Config{
//	       ClientID:     "... your application Client ID ...",
//	       ClientSecret: "... your application Client Secret ...",
//			  // ...
//	     }
//	     tc := oc.Client(ctx, "... your user Oauth Token ...")
//
//	     config := gphotos.Config{Client: tc}
//	     client, err := gphotos.New(config)
//	     ...
//	   }
//
// Note that hen using an authenticated client, all calls will include the specified OAuth 2.0 token. Therefore,
// authenticated clients should never be shared between different users.
//
// # Limitations
//
// Google Photos API imposes some limitations, please read them all at:
// https://github.com/gphotosuploader/google-photos-api-client-go/
//
// [this guide]: https://developers.google.com/photos/library/guides/get-started
// [Google Photos API]: https://developers.google.com/photos
package gphotos // import "github.com/gphotosuploader/google-photos-api-client-go/v3"
