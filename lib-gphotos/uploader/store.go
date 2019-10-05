package uploader

type UploadSessionStore interface {
	Get(fingerprint string) []byte
	Set(fingerprint string, url []byte)
	Delete(fingerprint string)
}
