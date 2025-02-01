package cloudinary

import (
	"net/http"
	"path/filepath"
	"strings"
	"trackonomy/internal/response"

	"github.com/gin-gonic/gin"
)

// CloudinaryController handles the routes for uploading images.
type CloudinaryController struct {
	service Service
}

func NewCloudinaryController(service Service) *CloudinaryController {
	return &CloudinaryController{service: service}
}

// UploadImage handles the image upload to Cloudinary.
func (cc *CloudinaryController) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Failed to retrieve file from request", err.Error())
		return
	}

	// Basic file-size validation (example: limit to ~2MB)
	if file.Size > 2*1024*1024 {
		response.BadRequest(c, "File size too large (max 2MB)", nil)
		return
	}

	// Simple extension check (you could also check MIME type, etc.)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExt[ext] {
		response.BadRequest(c, "Unsupported file format. Allowed: jpg, jpeg, png", nil)
		return
	}

	// Upload to Cloudinary
	url, err := cc.service.UploadImage(file)
	if err != nil {
		response.InternalServerError(c, "Failed to upload image", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "File uploaded successfully", gin.H{"url": url})
}
