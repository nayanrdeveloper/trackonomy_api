// internal/upload/validator.go
package upload

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// ValidateFile checks both file extension and file size in a single call.
func ValidateFile(
	fileHeader *multipart.FileHeader,
	allowedExtensions []string,
	maxSize int64,
) error {
	// Check the file type
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	var validExt bool
	for _, ext := range allowedExtensions {
		if fileExt == ext {
			validExt = true
			break
		}
	}
	if !validExt {
		return errors.New("invalid file type")
	}

	// Check the file size
	if fileHeader.Size > maxSize {
		return errors.New("file size exceeds the maximum limit")
	}

	return nil
}
