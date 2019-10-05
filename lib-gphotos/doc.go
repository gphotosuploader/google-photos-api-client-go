/*
Package gphotos provides a client for using the Google Photos API.
Wraps the Google Photos package provided by Google, and now removed:
https://github.com/gphotosuploader/googlemirror.

Usage:
	import "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"

Construct a new Google Photos client, then use the various services on the client to
access different parts of the Google Photos API. For example:

    // httpClient is an authenticated http.Client. See Authentication below.
	client := gphotos.NewClient(httpClient)
    // get or create a Photos Album with the specified name.
	album, err := GetOrCreateAlbumByName("my-new-album")
	// upload an specified file to an existent Photos Album.
    _, err := client.AddMediaItem(ctx, path, albumID)


NOTE: Using the https://godoc.org/context package, one can easily
pass cancellation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then context.Background()
can be used as a starting point.

Authentication
The gphotos library does not directly handle authentication. Instead, when
creating a new client, pass an http.Client that can handle authentication for
you. The easiest and recommended way to do this is using the golang.org/x/oauth2
library, but you can always use any other library that provides an http.Client.
Access to the API requires OAuth client credentials from a Google developers
project. This project must have the Library API enabled as described in
https://developers.google.com/photos/library/guides/get-started.

	import "golang.org/x/oauth2"
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

Note that when using an authenticated Client, all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.
See the oauth2 docs for complete instructions on using that library.

Rate Limiting
Google Photos imposes a rate limit on all API clients. The quota limit for
requests to the Library API is 10,000 requests per project per day. The quota
limit for requests to access media bytes (by loading a photo or video from a base
URL) is 75,000 requests per project per day.

Photo storage and quality
All media items uploaded to Google Photos using the API are stored in full
resolution at original quality (https://support.google.com/photos/answer/6220791).
They count toward the user’s storage.
*/
package gphotos
