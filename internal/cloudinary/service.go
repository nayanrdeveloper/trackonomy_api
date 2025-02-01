package cloudinary

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Service defines the interface for Cloudinary operations.
type Service interface {
	UploadImage(file *multipart.FileHeader) (string, error)
}

type service struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryService initializes a Cloudinary client service.
func NewCloudinaryService(cld *cloudinary.Cloudinary) Service {
	return &service{cld: cld}
}

// UploadImage uploads the provided file to Cloudinary and returns the secure URL.
func (s *service) UploadImage(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	uploadResult, err := s.cld.Upload.Upload(context.TODO(), src, uploader.UploadParams{
		Folder: "expense", // Adjust folder name as you like
	})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}
