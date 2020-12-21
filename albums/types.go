package albums

// Album represents a Google Photos album.
// See: https://developers.google.com/photos/library/reference/rest/v1/albums
type Album struct {
	ID                    string
	Title                 string
	ProductURL            string
	IsWriteable           bool
	MediaItemsCount       string
	CoverPhotoBaseURL     string
	CoverPhotoMediaItemID string
}
