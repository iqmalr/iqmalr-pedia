package services

import (
	"fmt"
	"mime/multipart"

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

	if err := clients.ValidateImageFile(file); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	imageURL, err := s.cloudinaryClient.UploadProductImage(src, productID)
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
		publicID := clients.ExtractPublicIDFromURL(imageURL)
		if publicID != "" {
			_ = s.cloudinaryClient.DeleteImage(publicID)
		}
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
		publicID := clients.ExtractPublicIDFromURL(image.ImageURL)
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
