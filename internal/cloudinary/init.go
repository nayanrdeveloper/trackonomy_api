package cloudinary

import (
	"trackonomy/config"

	"github.com/cloudinary/cloudinary-go/v2"
)

// InitClient initializes and returns a Cloudinary client using configuration from cfg.
func InitClient(cfg *config.Config) (*cloudinary.Cloudinary, error) {
	return cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)
}
