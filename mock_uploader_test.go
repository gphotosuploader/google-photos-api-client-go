package gphotos

import (
	"context"
)

// MockedUploader mocks an uploading service.
type MockedUploader struct {
	UploadFileFn func(ctx context.Context, filepath string) (string, error)
}

// Upload invokes the mock implementation and marks the function as invoked.
func (u MockedUploader) UploadFile(ctx context.Context, filepath string) (string, error) {
	return u.UploadFileFn(ctx, filepath)
}
