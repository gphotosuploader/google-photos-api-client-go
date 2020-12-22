package media_items

// MediaItem represents of a media item (such as a photo or video) in Google Photos.
// See: https://developers.google.com/photos/library/reference/rest/v1/mediaItems
type MediaItem struct {
	ID            string
	Description   string
	ProductURL    string
	BaseURL       string
	MimeType      string
	MediaMetadata MediaMetadata
	Filename      string
}

// MediaMetadata represents metadata for a media item.
// See: https://developers.google.com/photos/library/reference/rest/v1/mediaItems
type MediaMetadata struct {
	CreationTime string
	Width        string
	Height       string
}

// SimpleMediaItem represents a simple media item to be created in Google Photos via an upload token.
// See: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchCreate
type SimpleMediaItem struct {
	UploadToken string
	FileName    string
}
