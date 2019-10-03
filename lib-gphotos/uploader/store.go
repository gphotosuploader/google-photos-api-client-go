package uploader

type UploadSessionStore interface {
	Get(fingerprint string) (string, error)
	Set(fingerprint, url string) error
	Delete(fingerprint string) error
	Close()
}
