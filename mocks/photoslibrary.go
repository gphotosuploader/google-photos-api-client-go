package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// MockedGooglePhotosService mocks the Google Photos service.
type MockedGooglePhotosService struct {
	server  *httptest.Server
	baseURL string
}

const (
	// OK is returned on success.
	// @see: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
	grpcOKCode = 0
	// Unknown error. An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	//
	// The gRPC framework will generate this error code in the above two
	// mentioned cases.
	// @see: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
	grpcUnknownCode = 2
)

var (
	// AvailableAlbums is the albums collection.
	AvailableAlbums = []*photoslibrary.Album{
		{
			Id:         "fooId",
			ProductUrl: "fooProductUrl",
			Title:      "fooTitle",
		},
		{
			Id:         "barId",
			ProductUrl: "barProductUrl",
			Title:      "barTitle",
		},
		{
			Id:         "bazId",
			ProductUrl: "bazProductUrl",
			Title:      "bazTitle",
		},
	}

	// ShouldFailAlbum is an album that will make the API fail.
	ShouldFailAlbum = &photoslibrary.Album{
		Id:    "should-fail",
		Title: "should-fail",
	}

	// ShouldReturnPaginatedAlbum is an album that will contain a token to paginate over it.
	ShouldReturnPaginatedAlbum = &photoslibrary.Album{
		Id:    "should-return-paginated-album",
		Title: "should-return-paginated-album",
	}
)

// NewMockedGooglePhotosService returns a mocked Google Photos service.
func NewMockedGooglePhotosService() *MockedGooglePhotosService {
	ms := &MockedGooglePhotosService{}
	router := mux.NewRouter()
	// Albums methods
	router.HandleFunc("/v1/albums", ms.albumsList).Methods("GET")
	router.HandleFunc("/v1/albums", ms.albumsCreate).Methods("POST")
	router.HandleFunc("/v1/albums/{albumId}", ms.albumsGet).Methods("GET")
	router.HandleFunc("/v1/albums/{albumId}:batchAddMediaItems", ms.albumsBatchAddMediaItems).Methods("POST")
	// MediaItems methods
	router.HandleFunc("/v1/mediaItems:batchCreate", ms.mediaItemsBatchCreate).Methods("POST")
	router.HandleFunc("/v1/mediaItems/{mediaItemId}", ms.mediaItemsGet).Methods("GET")
	router.HandleFunc("/v1/mediaItems:search", ms.mediaItemsSearch).Methods("POST")

	ms.server = httptest.NewServer(router)
	ms.baseURL = ms.server.URL
	return ms
}

// Close closes the HTTP server.
func (ms MockedGooglePhotosService) Close() {
	ms.server.Close()
}

// URL returns the HTTP server url.
func (ms MockedGooglePhotosService) URL() string {
	return ms.baseURL
}

// albumsCreate implements 'albums.create' method .
// - Album creation with title == ShouldFailAlbum.Title will response http.StatusInternalServerError.
// - Any other case will response http.StatusOK.
//
// "flatPath": "v1/albums",
// "httpMethod": "POST",
func (ms MockedGooglePhotosService) albumsCreate(w http.ResponseWriter, r *http.Request) {
	var req photoslibrary.CreateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ShouldFailAlbum.Title == req.Album.Title {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	album := photoslibrary.Album{
		Id:         req.Album.Title + "Id",
		Title:      req.Album.Title + "Title",
		ProductUrl: req.Album.Title + "ProductUrl",
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// albumsGet implements 'albums.get' method.
// - Album with Id == ShouldFailAlbum.Title will response http.StatusInternalServerError.
// - Album with Id in AvailableAlbums will response http.StatusOK.
// - Any other case will response http.StatusNotFound.
//
// "flatPath": "v1/albums/{albumsId}",
// "httpMethod": "GET",
func (ms MockedGooglePhotosService) albumsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	if albumId == ShouldFailAlbum.Id {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	album, found := findAlbumById(albumId)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := album
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// albumsList implements 'albums.list' method.
//
// "flatPath": "v1/albums",
// "httpMethod": "GET",
func (ms MockedGooglePhotosService) albumsList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	res := photoslibrary.ListAlbumsResponse{
		Albums: AvailableAlbums,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// albumsBatchAddMediaItems implements 'albums.batchAddMediaItems' method.
//
// "flatPath": "v1/albums/{albumsId}:batchAddMediaItems",
// "httpMethod": "POST",
func (ms MockedGooglePhotosService) albumsBatchAddMediaItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	if _, found := findAlbumById(albumId); !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if ShouldFailAlbum.Id == albumId {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req photoslibrary.AlbumBatchAddMediaItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, mi := range req.MediaItemIds {
		if ShouldMakeAPIFailMediaItem == mi {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// findAlbumById returns if AvailableAlbums has an album with the specified Id.
func findAlbumById(albumId string) (*photoslibrary.Album, bool) {
	for _, a := range AvailableAlbums {
		if albumId == a.Id {
			return a, true
		}
	}
	return &photoslibrary.Album{}, false
}

var (
	// AvailableMediaItems is the number of media items in the fake collection.
	AvailableMediaItems = 30
	// DefaultPageSize is the number of items by page in results.
	DefaultPageSize = 10

	// ShouldMakeAPIFailMediaItem will make API fail.
	ShouldMakeAPIFailMediaItem = "should-make-API-fail"
	// ShouldReturnEmptyMediaItem will return an empty media item
	ShouldReturnEmptyMediaItem = "should-return-empty-media-item"
)

// albumsBatchRemoveMediaItems implements 'mediaItems.batchCreate' method.
//
// "flatPath": "v1/mediaItems:batchCreate",
// "httpMethod": "POST",
func (ms MockedGooglePhotosService) mediaItemsBatchCreate(w http.ResponseWriter, r *http.Request) {
	var req photoslibrary.BatchCreateMediaItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ShouldFailAlbum.Id == req.AlbumId {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newMediaItems := make([]*photoslibrary.NewMediaItemResult, len(req.NewMediaItems))
	for i, item := range req.NewMediaItems {
		if ShouldMakeAPIFailMediaItem == item.SimpleMediaItem.UploadToken {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if ShouldReturnEmptyMediaItem == item.SimpleMediaItem.UploadToken {
			newMediaItems[i] = &photoslibrary.NewMediaItemResult{
				Status: &photoslibrary.Status{Code: grpcUnknownCode},
			}
		} else {
			newMediaItems[i] = &photoslibrary.NewMediaItemResult{
				Status: &photoslibrary.Status{Code: grpcOKCode},
				MediaItem: &photoslibrary.MediaItem{
					BaseUrl:     item.SimpleMediaItem.UploadToken + "BaseUrl",
					Description: item.SimpleMediaItem.UploadToken + "Description",
					Filename:    item.SimpleMediaItem.UploadToken + "Filename",
					Id:          item.SimpleMediaItem.UploadToken + "Id",
					ProductUrl:  item.SimpleMediaItem.UploadToken + "ProductUrl",
					MediaMetadata: &photoslibrary.MediaMetadata{
						CreationTime: "2014-10-02T15:01:23.045123456Z",
						Height:       800,
						Width:        600,
					},
				},
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	res := photoslibrary.BatchCreateMediaItemsResponse{
		NewMediaItemResults: newMediaItems,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// mediaItemsGet implements 'mediaItems.get' method.
//
// "flatPath": "v1/mediaItems/{mediaItemId}",
// "httpMethod": "GET",
func (ms MockedGooglePhotosService) mediaItemsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaItemId := vars["mediaItemId"]

	if ShouldMakeAPIFailMediaItem == mediaItemId {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mediaItem, found := findMediaItemById(mediaItemId)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := mediaItem
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// mediaItemsSearch implements 'mediaItems.search' method.
//
// "flatPath": "v1/mediaItems:search",
// "httpMethod": "POST",
func (ms MockedGooglePhotosService) mediaItemsSearch(w http.ResponseWriter, r *http.Request) {
	var req photoslibrary.SearchMediaItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ShouldFailAlbum.Id == req.AlbumId {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusOK)
	res := photoslibrary.SearchMediaItemsResponse{
		MediaItems: createFakeMediaItems(AvailableMediaItems),
	}

	if ShouldReturnPaginatedAlbum.Id == req.AlbumId {
		if req.PageSize == 0 {
			req.PageSize = int64(DefaultPageSize)
		}

		totalItems := int64(len(res.MediaItems))
		pageNumber := getPageNumberFromToken(req.PageToken)
		pageStartsAt := pageNumber * req.PageSize
		pageEndsAt := pageStartsAt + req.PageSize

		res.MediaItems = res.MediaItems[pageStartsAt:pageEndsAt]
		nextPageNumber := pageNumber + 1
		if totalItems > pageEndsAt {
			res.NextPageToken = fmt.Sprintf("next-page-token-%d", nextPageNumber)
		}
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getPageNumberFromToken returns the number of page. Tokens are in the form 'next-page-token-<NUMBER>'.
func getPageNumberFromToken(token string) int64 {
	i := strings.Index(token, "next-page-token-")
	if i < 0 {
		return 0
	}
	if pageNumber, err := strconv.Atoi(token[i+len("next-page-token-"):]); err == nil {
		return int64(pageNumber)
	}
	return 0
}

// findMediaItemById returns if fake mediaItems collection has a media item with the specified Id.
func findMediaItemById(mediaItemId string) (*photoslibrary.MediaItem, bool) {
	for _, a := range createFakeMediaItems(AvailableMediaItems) {
		if mediaItemId == a.Id {
			return a, true
		}
	}
	return &photoslibrary.MediaItem{}, false
}

// createFakeMediaItems returns a collection of MediaItems with the specified number of it.
func createFakeMediaItems(numberOfItems int) []*photoslibrary.MediaItem {

	mediaItemsResult := make([]*photoslibrary.MediaItem, numberOfItems)

	for i := 0; i < numberOfItems; i++ {
		mediaItemsResult[i] = &photoslibrary.MediaItem{
			Id:          fmt.Sprintf("fooId-%d", i),
			Description: fmt.Sprintf("fooDescription-%d", i),
			ProductUrl:  fmt.Sprintf("fooProductUrl-%d", i),
			BaseUrl:     fmt.Sprintf("fooBaseUrl-%d", i),
			Filename:    fmt.Sprintf("fooFilename-%d", i),
			MediaMetadata: &photoslibrary.MediaMetadata{
				CreationTime: "2014-10-02T15:01:23.045123456Z",
				Height:       800,
				Width:        600,
			},
		}
	}

	return mediaItemsResult
}
