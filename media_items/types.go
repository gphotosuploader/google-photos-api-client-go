package media_items

type SimpleMediaItem struct {
	UploadToken string
	FileName    string
}

type MediaItem struct {
	ID            string
	Description   string
	ProductURL    string
	BaseURL       string
	MimeType      string
	MediaMetadata MediaMetadata
	Filename      string
}

type MediaMetadata struct {
	CreationTime string
	Width        string
	Height       string
}
