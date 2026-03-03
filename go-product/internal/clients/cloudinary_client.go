package clients

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/iqmalr-pedia/go-product/internal/config"
)

const (
	RootFolder     = "iqmalr-pedia"
	MaxFileSize    = 2 * 1024 * 1024 // 2MB
)

var AllowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

type ImageCategory string

const (
	CategoryProducts ImageCategory = "products"
)

type CloudinaryClient interface {
	UploadProductImage(file multipart.File, productID uint) (string, error)
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

func (c *cloudinaryClient) UploadProductImage(file multipart.File, productID uint) (string, error) {
	ctx := context.Background()
	publicID := BuildPublicID(CategoryProducts, "product", productID)
	folder := fmt.Sprintf("%s/%s", RootFolder, CategoryProducts)

	result, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       publicID,
		Overwrite:      boolPtr(false),
		UniqueFilename: boolPtr(false),
		Transformation: "w_800,h_800,c_fill,g_auto,q_auto,f_webp",
	})
	if err != nil {
		log.Printf("Cloudinary upload error for product %d: %v", productID, err)
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return result.SecureURL, nil
}

func (c *cloudinaryClient) DeleteImage(publicID string) error {
	ctx := context.Background()

	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		log.Printf("Cloudinary delete error for %s: %v", publicID, err)
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

func BuildPublicID(category ImageCategory, entity string, id uint) string {
	return fmt.Sprintf("%s-%d-%d", entity, id, time.Now().Unix())
}

func ValidateImageFile(header *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedImageExts[ext] {
		return fmt.Errorf("file type '%s' is not allowed, use jpg, jpeg, png, or webp", ext)
	}
	if header.Size > MaxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum limit of 2MB", header.Size)
	}
	return nil
}

func ExtractPublicIDFromURL(url string) string {
	idx := strings.Index(url, "/upload/")
	if idx == -1 {
		return ""
	}
	path := url[idx+len("/upload/"):]

	if strings.HasPrefix(path, "v") {
		if slashIdx := strings.Index(path, "/"); slashIdx != -1 {
			path = path[slashIdx+1:]
		}
	}

	ext := filepath.Ext(path)
	if ext != "" {
		path = path[:len(path)-len(ext)]
	}

	return path
}

func boolPtr(b bool) *bool {
	return &b
}
