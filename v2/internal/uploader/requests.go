package uploader

import (
	"fmt"
	"net/http"
	"path"
)

// createRawUploadRequest returns a raw (non-resumable) upload request for Google Photos.
func createRawUploadRequest(url string, upload *Upload) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, upload.r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", path.Base(upload.name))
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	return req, nil

}

// createInitialResumableUploadRequest returns a starting resumable upload request for Google Photos.
// url is the Google Photos upload endpoint.
func createInitialResumableUploadRequest(url string, upload *Upload) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", path.Base(upload.name))
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Raw-Size", fmt.Sprintf("%d", upload.size))

	return req, nil
}

// createQueryOffsetRequest returns a query offset request for Google Photos.
// url is the unique URL that must be used to complete the upload through all of the remaining requests.
func createQueryOffsetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "query")

	return req, nil
}

// createResumeUploadRequest returns a resumable upload request to continue an started upload for Google Photos.
// url is the unique URL that must be used to complete the upload through all of the remaining requests.
func createResumeUploadRequest(url string, upload *Upload) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, upload.r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Length", fmt.Sprintf("%d", upload.size-upload.sent))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", upload.sent))

	return req, nil
}
