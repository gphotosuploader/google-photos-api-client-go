package uploader

type UploadSessionStore interface {
	Get(fingerprint string) string
	Set(fingerprint, url string)
	Delete(fingerprint string)
}
