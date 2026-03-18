package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryService handles file uploads to Cloudinary.
type CloudinaryService struct {
	cloudName string
	apiKey    string
	apiSecret string
}

// NewCloudinaryService constructs a CloudinaryService with the provided credentials.
func NewCloudinaryService(cloudName, apiKey, apiSecret string) *CloudinaryService {
	return &CloudinaryService{
		cloudName: cloudName,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// UploadFile opens a multipart file and uploads it to Cloudinary under the given folder.
// Returns the secure HTTPS URL of the uploaded asset.
func (s *CloudinaryService) UploadFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("cloudinary: open file: %w", err)
	}
	defer file.Close()

	cld, err := cloudinary.NewFromParams(s.cloudName, s.apiKey, s.apiSecret)
	if err != nil {
		return "", fmt.Errorf("cloudinary: init client: %w", err)
	}

	result, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary: upload: %w", err)
	}

	return result.SecureURL, nil
}
