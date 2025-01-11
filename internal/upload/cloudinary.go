package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"trackonomy/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryService is an interface that defines methods for uploading files to Cloudinary.
type CloudinaryService interface {
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error)
}

type cloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryService initializes a Cloudinary client and returns a CloudinaryService.
func NewCloudinaryService(cfg *config.Config) (CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary client: %w", err)
	}
	return &cloudinaryService{cld: cld}, nil
}

// UploadFile uploads a file (image/pdf/etc.) to Cloudinary and returns its URL.
func (cs *cloudinaryService) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Prepare the public ID or let Cloudinary auto-generate one
	// Example: you might sanitize filename, or you can leave it for Cloudinary to handle
	fileName := fileHeader.Filename
	publicID := strings.TrimSuffix(fileName, ".pdf") // example for PDF, similarly for images

	uploadParams := uploader.UploadParams{
		Folder:   folder,   // e.g. "trackonomy/expenses"
		PublicID: publicID, // optional - set if you want custom naming
	}

	// Actually perform the upload
	uploadResult, err := cs.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return uploadResult.SecureURL, nil
}
