package mocks

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

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
	// an error-space that is not known in this address space. Also,
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	//
	// The gRPC framework will generate this error code in the above two
	// mentioned cases.
	// @see: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
	grpcUnknownCode = 2

	// maxMediaItemsPerPage is the maximum number of media items to request from the PhotosLibrary. Fewer media items
	// might be returned than the specified number.
	// See https://developers.google.com/photos/library/guides/list#pagination
	maxMediaItemsPerPage = 100

	// maxAlbumsPerPage is the maximum number of albums to request from the PhotosLibrary. Fewer albums
	// might be returned than the specified number.
	// See https://developers.google.com/photos/library/guides/list#pagination
	maxAlbumsPerPage = 50

	// AvailableMediaItems is the number of media items in the fake collection. It should be bigger than `maxItemsPerPage`.
	AvailableMediaItems = 150
	// AvailableAlbums is the number of media items in the fake collection. It should be bigger than `maxItemsPerPage`.
	AvailableAlbums = 75

	// ShouldMakeAPIFailMediaItem will make API fail.
	ShouldMakeAPIFailMediaItem = "should-make-API-fail"
	// ShouldReturnEmptyMediaItem will return an empty media item
	ShouldReturnEmptyMediaItem = "should-return-empty-media-item"
)

var (
	// ShouldFailAlbum is an album that will make the API fail.
	ShouldFailAlbum = &photoslibrary.Album{
		Id:    "should-fail",
		Title: "should-fail",
	}
)

// NewMockedGooglePhotosService returns a mocked Google Photos service.
func NewMockedGooglePhotosService() *MockedGooglePhotosService {
	ms := &MockedGooglePhotosService{}
	router := chi.NewRouter()
	// Albums methods
	router.Get("/v1/albums", ms.albumsList)
	router.Post("/v1/albums", ms.albumsCreate)
	router.Get("/v1/albums/{albumId}", ms.albumsGet)
	router.Post("/v1/albums/{albumId}:batchAddMediaItems", ms.albumsBatchAddMediaItems)
	// MediaItems methods
	router.Post("/v1/mediaItems:batchCreate", ms.mediaItemsBatchCreate)
	router.Get("/v1/mediaItems/{mediaItemId}", ms.mediaItemsGet)
	router.Post("/v1/mediaItems:search", ms.mediaItemsSearch)
	// Uploads methods
	router.Post("/v1/uploads", ms.uploads)
	router.Post("/v1/upload-session/started", ms.handleExistingUploadSession)
	router.Post("/v1/upload-session/upload-success", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		_, _ = w.Write([]byte("apiToken"))
	})

	ms.server = httptest.NewServer(router)
	ms.baseURL = ms.server.URL
	return ms
}

// Close closes the HTTP server.
func (ms *MockedGooglePhotosService) Close() {
	ms.server.Close()
}

// URL returns the HTTP server url.
func (ms *MockedGooglePhotosService) URL() string {
	return ms.baseURL
}

// albumsCreate implements 'albums.create' method .
// - Album creation with title == ShouldFailAlbum.Title will respond http.StatusInternalServerError.
// - Any other case will respond http.StatusOK.
//
// "flatPath": "v1/albums",
// "httpMethod": "POST",
func (ms *MockedGooglePhotosService) albumsCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
// - Album with Id == ShouldFailAlbum.Id will respond http.StatusInternalServerError.
// - Album with Id in AvailableAlbums will response http.StatusOK.
// - Any other case will respond http.StatusNotFound.
//
// "flatPath": "v1/albums/{albumsId}",
// "httpMethod": "GET",
func (ms *MockedGooglePhotosService) albumsGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	albumId := chi.URLParam(r, "albumId")

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
func (ms *MockedGooglePhotosService) albumsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	albums := createFakeAlbums(AvailableAlbums)

	pageSize, pageToken := ms.paginationOptions(r)
	p := newAlbumsPaginator(pageSize, albums)
	items, nextPageToken := p.page(pageToken)

	w.WriteHeader(http.StatusOK)
	res := photoslibrary.ListAlbumsResponse{
		Albums:        items,
		NextPageToken: nextPageToken,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ms *MockedGooglePhotosService) paginationOptions(r *http.Request) (pageSize int64, pageToken string) {
	pt := r.URL.Query().Get("pageToken")
	ps, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		return maxAlbumsPerPage, pt
	}
	return int64(ps), pt
}

type albumsPaginator struct {
	items    []*photoslibrary.Album
	limit    int64
	pageSize int64
}

func newAlbumsPaginator(pageSize int64, items []*photoslibrary.Album) *albumsPaginator {
	if pageSize < 1 || pageSize > maxAlbumsPerPage {
		pageSize = maxAlbumsPerPage
	}

	return &albumsPaginator{
		limit:    int64(len(items)),
		items:    items,
		pageSize: pageSize,
	}
}

func (p *albumsPaginator) page(pageToken string) (results []*photoslibrary.Album, nextPageToken string) {
	maxPages := int64(len(p.items)) / p.pageSize
	pageNumber := getPageNumberFromToken(pageToken, maxPages)
	pageStartAt := p.pageSize * pageNumber
	pageEndsAt := pageStartAt + p.pageSize

	if pageEndsAt >= p.limit {
		return p.items[pageStartAt:], ""
	}

	nextPageToken = fmt.Sprintf("next-page-token-%d", pageNumber+1)
	return p.items[pageStartAt:pageEndsAt], nextPageToken
}

// getPageNumberFromToken returns the number of page.
// Tokens are in the form 'next-page-token-<NUMBER>'.
func getPageNumberFromToken(token string, maxPages int64) int64 {
	i := strings.Index(token, "next-page-token-")
	if i < 0 {
		return 0
	}
	pageNumber, err := strconv.Atoi(token[i+len("next-page-token-"):])
	if err != nil || int64(pageNumber) > maxPages {
		return 0
	}
	return int64(pageNumber)
}

func createFakeAlbums(numberOfItems int) []*photoslibrary.Album {
	albumsResult := make([]*photoslibrary.Album, numberOfItems)
	for i := 0; i < numberOfItems; i++ {
		albumsResult[i] = &photoslibrary.Album{
			Id:         fmt.Sprintf("fooId-%d", i),
			ProductUrl: fmt.Sprintf("fooProductUrl-%d", i),
			Title:      fmt.Sprintf("fooTitle-%d", i),
		}
	}
	return albumsResult
}

// albumsBatchAddMediaItems implements 'albums.batchAddMediaItems' method.
//
// "flatPath": "v1/albums/{albumsId}:batchAddMediaItems",
// "httpMethod": "POST",
func (ms *MockedGooglePhotosService) albumsBatchAddMediaItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	albumId := chi.URLParam(r, "albumId")

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
	for _, a := range createFakeAlbums(AvailableAlbums) {
		if albumId == a.Id {
			return a, true
		}
	}
	return &photoslibrary.Album{}, false
}

// albumsBatchRemoveMediaItems implements 'mediaItems.batchCreate' method.
//
// "flatPath": "v1/mediaItems:batchCreate",
// "httpMethod": "POST",
func (ms *MockedGooglePhotosService) mediaItemsBatchCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
func (ms *MockedGooglePhotosService) mediaItemsGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mediaItemId := chi.URLParam(r, "mediaItemId")

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
func (ms *MockedGooglePhotosService) mediaItemsSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req photoslibrary.SearchMediaItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ShouldFailAlbum.Id == req.AlbumId {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mediaItems := createFakeMediaItems(AvailableMediaItems)
	p := newMediaItemsPaginator(req.PageSize, mediaItems)
	items, nextPageToken := p.page(req.PageToken)
	res := photoslibrary.SearchMediaItemsResponse{
		MediaItems:    items,
		NextPageToken: nextPageToken,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type mediaItemsPaginator struct {
	items    []*photoslibrary.MediaItem
	limit    int64
	pageSize int64
}

func newMediaItemsPaginator(pageSize int64, items []*photoslibrary.MediaItem) *mediaItemsPaginator {
	if pageSize < 1 || pageSize > maxMediaItemsPerPage {
		pageSize = maxMediaItemsPerPage
	}

	return &mediaItemsPaginator{
		limit:    int64(len(items)),
		items:    items,
		pageSize: pageSize,
	}
}

func (p *mediaItemsPaginator) page(pageToken string) (results []*photoslibrary.MediaItem, nextPageToken string) {
	maxPages := int64(len(p.items)) / p.pageSize
	pageNumber := getPageNumberFromToken(pageToken, maxPages)
	pageStartAt := p.pageSize * pageNumber
	pageEndsAt := pageStartAt + p.pageSize

	if pageEndsAt >= p.limit {
		return p.items[pageStartAt:], ""
	}

	nextPageToken = fmt.Sprintf("next-page-token-%d", pageNumber+1)
	return p.items[pageStartAt:pageEndsAt], nextPageToken
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
func createFakeMediaItems(numberOfItems int64) []*photoslibrary.MediaItem {
	mediaItemsResult := make([]*photoslibrary.MediaItem, numberOfItems)
	for i := int64(0); i < numberOfItems; i++ {
		mediaItemsResult[i] = &photoslibrary.MediaItem{
			BaseUrl:     fmt.Sprintf("fooBaseUrl-%d", i),
			Description: fmt.Sprintf("fooDescription-%d", i),
			Filename:    fmt.Sprintf("fooFilename-%d", i),
			Id:          fmt.Sprintf("fooId-%d", i),
			ProductUrl:  fmt.Sprintf("fooProductUrl-%d", i),
			MediaMetadata: &photoslibrary.MediaMetadata{
				CreationTime: "2014-10-02T15:01:23.045123456Z",
				Height:       800,
				Width:        600,
			},
		}
	}
	return mediaItemsResult
}

func (ms *MockedGooglePhotosService) uploads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if "resumable" == r.Header.Get("X-Goog-Upload-Protocol") {
		ms.handleResumableUpload(w, r)
	}

	ms.handleSimpleUpload(w, r)
}

func (ms *MockedGooglePhotosService) handleSimpleUpload(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("X-Goog-Upload-File-Name") {
	case "upload-failure":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		var bodyContent []byte
		bodyLength, err := r.Body.Read(bodyContent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		expectedLength, _ := strconv.Atoi(r.Header.Get("Content-Length"))
		if expectedLength != bodyLength {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte("apiToken"))
	}
}

func (ms *MockedGooglePhotosService) handleResumableUpload(w http.ResponseWriter, r *http.Request) {
	if "start" != r.Header.Get("X-Goog-Upload-Command") {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Header.Get("X-Goog-Upload-File-Name") {
	case "upload-should-fail":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.Header().Add("X-Goog-Upload-URL", ms.URL()+"/v1/upload-session/upload-success")
	}
}

func (ms *MockedGooglePhotosService) handleExistingUploadSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.Header.Get("X-Goog-Upload-Command") {
	case "query":
		w.Header().Add("X-Goog-Upload-Status", "active")
		w.Header().Add("X-Goog-Upload-Size-Received", "1000")
	case "upload, finalize":
		_, _ = w.Write([]byte("apiToken"))
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
