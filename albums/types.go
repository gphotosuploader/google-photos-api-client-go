package albums

import (
	"net/http"

	"github.com/duffpl/google-photos-api-client/albums"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
)

type Album = albums.Album
type Field = albums.Field
type ListOptions = albums.AlbumsListOptions

func defaultCache() *cache.CachitaCache {
	return cache.NewCachitaCache()
}

func defaultRepo(authenticatedClient *http.Client) albums.HttpAlbumsService {
	return albums.NewHttpAlbumsService(authenticatedClient)
}

