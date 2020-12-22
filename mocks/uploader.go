package mocks

import (
	"context"
)

// MockedUploader mocks an uploader.
type MockedUploader struct {
	UploadFileFn func(ctx context.Context, filePath string) (uploadToken string, err error)
}

// UploadFile invokes the mock implementation.
func (m MockedUploader) UploadFile(ctx context.Context, filePath string) (uploadToken string, err error) {
	return m.UploadFileFn(ctx, filePath)
}
