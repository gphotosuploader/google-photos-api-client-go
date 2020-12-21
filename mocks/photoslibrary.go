package mocks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// MockedGooglePhotosService mock the Google Photos service.
type MockedGooglePhotosService struct {
	server  *httptest.Server
	baseURL string
}

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

	ShouldFailAlbum = &photoslibrary.Album{
		Id:    "should-fail",
		Title: "should-fail",
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

func (ms MockedGooglePhotosService) Close() {
	ms.server.Close()
}

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
		if ShouldFailMediaItem == mi {
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
	// AvailableMediaItems is the media items collection.
	AvailableMediaItems = []*photoslibrary.MediaItem{
		{
			Id:          "fooId",
			Description: "fooDescription",
			ProductUrl:  "fooProductUrl",
			BaseUrl:     "fooBaseUrl",
			Filename:    "fooFilename",
			MediaMetadata: &photoslibrary.MediaMetadata{
				CreationTime: "2014-10-02T15:01:23.045123456Z",
				Height:       800,
				Width:        600,
			},
		},
		{
			Id:          "barId",
			Description: "barDescription",
			ProductUrl:  "barProductUrl",
			BaseUrl:     "barBaseUrl",
			Filename:    "barFilename",
			MediaMetadata: &photoslibrary.MediaMetadata{
				CreationTime: "2014-10-02T15:01:23.045123456Z",
				Height:       800,
				Width:        600,
			},
		},
		{
			Id:          "bazId",
			Description: "bazDescription",
			ProductUrl:  "bazProductUrl",
			BaseUrl:     "bazBaseUrl",
			Filename:    "bazFilename",
			MediaMetadata: &photoslibrary.MediaMetadata{
				CreationTime: "2014-10-02T15:01:23.045123456Z",
				Height:       800,
				Width:        600,
			},
		},
	}
	ShouldFailMediaItem = "should-fail"
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
		if ShouldFailMediaItem == item.SimpleMediaItem.UploadToken {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newMediaItems[i] = &photoslibrary.NewMediaItemResult{
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

	if ShouldFailMediaItem == mediaItemId {
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
		MediaItems: AvailableMediaItems,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// findMediaItemById returns if AvailableMediaItems has a media item with the specified Id.
func findMediaItemById(mediaItemId string) (*photoslibrary.MediaItem, bool) {
	for _, a := range AvailableMediaItems {
		if mediaItemId == a.Id {
			return a, true
		}
	}
	return &photoslibrary.MediaItem{}, false
}
