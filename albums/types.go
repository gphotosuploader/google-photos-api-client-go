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

// maxItemsPerPage is the maximum number of albums to ask the PhotosLibrary.
// Fewer albums might be returned than the specified number.
// See https://developers.google.com/photos/library/guides/list#pagination
const maxItemsPerPage = 50

// Options define the options that could be customized when listing albums.
type Options struct {
	// ExcludeNonAppCreatedData excludes albums that were not created by this app.
	// Defaults to false (all albums are returned).
	ExcludeNonAppCreatedData bool
}
