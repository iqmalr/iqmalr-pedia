package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/iqmalr-pedia/go-product/internal/clients"
	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/models"
	"github.com/iqmalr-pedia/go-product/internal/repositories"
)

type ImageService interface {
	UploadImage(productID uint, file *multipart.FileHeader, altText string, sortOrder int, isPrimary bool) (*response.ImageDetailResponse, error)
	GetImagesByProductID(productID uint) (*response.ImageListResponse, error)
	UpdateImage(productID, imageID uint, req *request.UpdateImageRequest) (*response.ImageDetailResponse, error)
	DeleteImage(productID, imageID uint) error
	SetPrimaryImage(productID, imageID uint) error
}

type imageService struct {
	imageRepo        repositories.ImageRepository
	productRepo      repositories.ProductRepository
	cloudinaryClient clients.CloudinaryClient
}

func NewImageService(imageRepo repositories.ImageRepository, productRepo repositories.ProductRepository, cloudinaryClient clients.CloudinaryClient) ImageService {
	return &imageService{
		imageRepo:        imageRepo,
		productRepo:      productRepo,
		cloudinaryClient: cloudinaryClient,
	}
}

func (s *imageService) UploadImage(productID uint, file *multipart.FileHeader, altText string, sortOrder int, isPrimary bool) (*response.ImageDetailResponse, error) {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("file type '%s' is not allowed, use jpg, jpeg, png, gif, or webp", ext)
	}

	if file.Size > 2*1024*1024 {
		return nil, fmt.Errorf("file size exceeds maximum limit of 2MB")
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	folder := fmt.Sprintf("products/%d", productID)
	publicID := fmt.Sprintf("product_%d_%d", productID, time.Now().UnixNano())

	imageURL, err := s.cloudinaryClient.UploadImage(src, folder, publicID)
	if err != nil {
		return nil, err
	}

	if isPrimary {
		_ = s.imageRepo.ClearPrimaryByProductID(productID)
	}

	image := &models.ProductImage{
		ProductID: productID,
		ImageURL:  imageURL,
		AltText:   altText,
		SortOrder: sortOrder,
		IsPrimary: isPrimary,
	}

	if err := s.imageRepo.Create(image); err != nil {
		// best-effort cleanup on Cloudinary
		fullPublicID := folder + "/" + publicID
		_ = s.cloudinaryClient.DeleteImage(fullPublicID)
		return nil, err
	}

	return s.toImageDetailResponse(image), nil
}

func (s *imageService) GetImagesByProductID(productID uint) (*response.ImageListResponse, error) {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	images, err := s.imageRepo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}

	var items []response.ImageDetailResponse
	for _, img := range images {
		items = append(items, *s.toImageDetailResponse(&img))
	}

	return &response.ImageListResponse{Data: items}, nil
}

func (s *imageService) UpdateImage(productID, imageID uint, req *request.UpdateImageRequest) (*response.ImageDetailResponse, error) {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	image, err := s.imageRepo.GetByID(imageID)
	if err != nil {
		return nil, err
	}

	if image.ProductID != productID {
		return nil, fmt.Errorf("image does not belong to this product")
	}

	if req.AltText != "" {
		image.AltText = req.AltText
	}

	if req.SortOrder != nil {
		image.SortOrder = *req.SortOrder
	}

	if req.IsPrimary != nil && *req.IsPrimary {
		_ = s.imageRepo.ClearPrimaryByProductID(productID)
		image.IsPrimary = true
	} else if req.IsPrimary != nil {
		image.IsPrimary = false
	}

	if err := s.imageRepo.Update(image); err != nil {
		return nil, err
	}

	return s.toImageDetailResponse(image), nil
}

func (s *imageService) DeleteImage(productID, imageID uint) error {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	image, err := s.imageRepo.GetByID(imageID)
	if err != nil {
		return err
	}

	if image.ProductID != productID {
		return fmt.Errorf("image does not belong to this product")
	}

	if image.ImageURL != "" {
		publicID := extractPublicIDFromURL(image.ImageURL)
		if publicID != "" {
			_ = s.cloudinaryClient.DeleteImage(publicID)
		}
	}

	return s.imageRepo.Delete(imageID)
}

func (s *imageService) SetPrimaryImage(productID, imageID uint) error {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	image, err := s.imageRepo.GetByID(imageID)
	if err != nil {
		return err
	}

	if image.ProductID != productID {
		return fmt.Errorf("image does not belong to this product")
	}

	if err := s.imageRepo.ClearPrimaryByProductID(productID); err != nil {
		return err
	}

	image.IsPrimary = true
	return s.imageRepo.Update(image)
}

func (s *imageService) toImageDetailResponse(img *models.ProductImage) *response.ImageDetailResponse {
	return &response.ImageDetailResponse{
		ID:        img.ID,
		ProductID: img.ProductID,
		ImageURL:  img.ImageURL,
		AltText:   img.AltText,
		SortOrder: img.SortOrder,
		IsPrimary: img.IsPrimary,
		CreatedAt: img.CreatedAt,
		UpdatedAt: img.UpdatedAt,
	}
}

// extractPublicIDFromURL extracts Cloudinary public_id from a secure URL.
// e.g. https://res.cloudinary.com/xxx/image/upload/v123/products/1/product_1_123.jpg
// returns "products/1/product_1_123"
func extractPublicIDFromURL(url string) string {
	idx := strings.Index(url, "/upload/")
	if idx == -1 {
		return ""
	}
	path := url[idx+len("/upload/"):]

	// skip version segment (v1234567890/)
	if strings.HasPrefix(path, "v") {
		if slashIdx := strings.Index(path, "/"); slashIdx != -1 {
			path = path[slashIdx+1:]
		}
	}

	// remove file extension
	ext := filepath.Ext(path)
	if ext != "" {
		path = path[:len(path)-len(ext)]
	}

	return path
}
