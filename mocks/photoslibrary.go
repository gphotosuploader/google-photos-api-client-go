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
	mux := mux.NewRouter()
	mux.HandleFunc("/v1/albums", ms.handleGetAlbumsList).Methods("GET")
	mux.HandleFunc("/v1/albums", ms.handlePostAlbumsCreate).Methods("POST")
	mux.HandleFunc("/v1/albums/{albumId}", ms.handleGetAlbumsGet).Methods("GET")

	ms.server = httptest.NewServer(mux)
	ms.baseURL = ms.server.URL
	return ms
}

func (ms MockedGooglePhotosService) Close() {
	ms.server.Close()
}

func (ms MockedGooglePhotosService) URL() string {
	return ms.baseURL
}

// handleCreateAlbum implements 'albums.create' method .
// - Album creation with title == ShouldFailAlbum.Title will response http.StatusInternalServerError.
// - Any other case will response http.StatusCreated.
//
// "flatPath": "v1/albums",
// "httpMethod": "POST",
func (ms MockedGooglePhotosService) handlePostAlbumsCreate(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// handleGetAlbum implements 'albums.get' method.
// - Album with Id == ShouldFailAlbum.Title will response http.StatusInternalServerError.
// - Album with Id in AvailableAlbums will response http.StatusOK.
// - Any other case will response http.StatusNotFound.
//
// "flatPath": "v1/albums/{albumsId}",
// "httpMethod": "GET",
func (ms MockedGooglePhotosService) handleGetAlbumsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	if albumId == ShouldFailAlbum.Id {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	album, found := findById(albumId)
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

// handleListAlbum implements 'albums.list' method.
//
// "flatPath": "v1/albums",
// "httpMethod": "GET",
func (ms MockedGooglePhotosService) handleGetAlbumsList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	res := photoslibrary.ListAlbumsResponse{
		Albums: AvailableAlbums,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// findById returns if AvailableAlbums has an album with the specified Id.
func findById(albumId string) (*photoslibrary.Album, bool) {
	for _, a := range AvailableAlbums {
		if albumId == a.Id {
			return a, true
		}
	}
	return &photoslibrary.Album{}, false
}
