package clients

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/iqmalr-pedia/go-product/internal/config"
)

type CloudinaryClient interface {
	UploadImage(file multipart.File, folder string, publicID string) (string, error)
	DeleteImage(publicID string) error
}

type cloudinaryClient struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryClient() (CloudinaryClient, error) {
	cld, err := cloudinary.NewFromParams(
		config.AppConfig.CloudinaryCloudName,
		config.AppConfig.CloudinaryAPIKey,
		config.AppConfig.CloudinaryAPISecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &cloudinaryClient{cld: cld}, nil
}

func (c *cloudinaryClient) UploadImage(file multipart.File, folder string, publicID string) (string, error) {
	ctx := context.Background()

	uploadResult, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   folder,
		PublicID: publicID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}

func (c *cloudinaryClient) DeleteImage(publicID string) error {
	ctx := context.Background()

	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %w", err)
	}

	return nil
}
