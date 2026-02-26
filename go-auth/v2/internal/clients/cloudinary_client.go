package clients

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	RootFolder  = "iqmalr-pedia"
	MaxFileSize = 2 * 1024 * 1024 // 2MB
)

var AllowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

type CloudinaryClient interface {
	UploadAvatar(file multipart.File, userID uint) (string, error)
	DeleteImage(publicID string) error
}

type cloudinaryClient struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryClient() (CloudinaryClient, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &cloudinaryClient{cld: cld}, nil
}

func (c *cloudinaryClient) UploadAvatar(file multipart.File, userID uint) (string, error) {
	ctx := context.Background()
	publicID := fmt.Sprintf("user-%d-%d", userID, time.Now().Unix())
	folder := fmt.Sprintf("%s/avatars", RootFolder)

	result, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		PublicID:       publicID,
		Overwrite:      boolPtr(false),
		UniqueFilename: boolPtr(false),
		Transformation: "w_500,h_500,c_fill,g_auto,q_auto,f_webp",
	})
	if err != nil {
		log.Printf("Cloudinary upload error for user %d: %v", userID, err)
		return "", fmt.Errorf("failed to upload avatar: %w", err)
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
