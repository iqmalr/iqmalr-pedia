package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/iqmalr-pedia/go-product/internal/clients"
	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/models"
	"github.com/iqmalr-pedia/go-product/internal/repositories"
	"github.com/iqmalr-pedia/go-product/internal/utils"
	_ "gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(req *request.CreateProductRequest) (*response.ProductResponse, error)
	GetProductByID(id uint) (*response.ProductDetailResponse, error)
	GetProductBySlug(slug string) (*response.ProductDetailResponse, error)
	UpdateProduct(id uint, req *request.UpdateProductRequest) (*response.ProductResponse, error)
	DeleteProduct(id uint) error
	ListProducts(req *request.ListProductsRequest) (*response.ProductListResponse, error)
	UpdateProductStatus(id uint, req *request.UpdateProductStatusRequest) (*response.ProductStatusResponse, error)
	PublishProduct(id uint) (*response.ProductPublishResponse, error)
	UnpublishProduct(id uint) (*response.ProductPublishResponse, error)
}

type VendorRepository interface {
	GetByID(id uint) (*models.Vendor, error)
}

type productService struct {
	productRepo   repositories.ProductRepository
	categoryRepo  repositories.CategoryRepository
	vendorService VendorService
}

func NewProductService(productRepo repositories.ProductRepository, categoryRepo repositories.CategoryRepository, vendorService VendorService) ProductService {
	return &productService{
		productRepo:   productRepo,
		categoryRepo:  categoryRepo,
		vendorService: vendorService,
	}
}

func (s *productService) CreateProduct(req *request.CreateProductRequest) (*response.ProductResponse, error) {
	_, err := s.vendorService.GetVendorByID(req.VendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vendor information: %w", err)
	}

	trackInventory := true
	if req.TrackInventory != nil {
		trackInventory = *req.TrackInventory
	}

	allowBackorder := false
	if req.AllowBackorder != nil {
		allowBackorder = *req.AllowBackorder
	}

	isFeatured := false
	if req.IsFeatured != nil {
		isFeatured = *req.IsFeatured
	}

	status := "draft"
	if req.Status != "" {
		status = req.Status
	}

	product := &models.Product{
		VendorID:          req.VendorID,
		Name:              req.Name,
		Description:       req.Description,
		ShortDescription:  req.ShortDescription,
		Price:             req.Price,
		CompareAtPrice:    req.CompareAtPrice,
		CostPerItem:       req.CostPerItem,
		Stock:             req.Stock,
		LowStockThreshold: req.LowStockThreshold,
		TrackInventory:    trackInventory,
		AllowBackorder:    allowBackorder,
		Weight:            req.Weight,
		Length:            req.Length,
		Width:             req.Width,
		Height:            req.Height,
		MetaTitle:         req.MetaTitle,
		MetaDescription:   req.MetaDescription,
		MetaKeywords:      req.MetaKeywords,
		Status:            status,
		IsFeatured:        isFeatured,
	}

	if len(req.CategoryIDs) > 0 {
		var categories []models.Category
		for _, categoryID := range req.CategoryIDs {
			category, err := s.categoryRepo.GetByID(categoryID)
			if err != nil {
				return nil, fmt.Errorf("category with ID %d not found", categoryID)
			}
			categories = append(categories, *category)
		}
		product.Categories = categories
	}

	err = s.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productService) GetProductByID(id uint) (*response.ProductDetailResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	//return s.toProductDetailResponse(product), nil
	return s.toProductDetailResponse(product)
}

func (s *productService) GetProductBySlug(slug string) (*response.ProductDetailResponse, error) {
	product, err := s.productRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	//return s.toProductDetailResponse(product), nil
	return s.toProductDetailResponse(product)
}

func (s *productService) UpdateProduct(id uint, req *request.UpdateProductRequest) (*response.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" && req.Name != product.Name {
		product.Name = req.Name
		shortID := strings.Split(product.UUID.String(), "-")[0]
		product.Slug = utils.GenerateSlug(product.Name, shortID)
	}

	if req.Description != "" {
		product.Description = req.Description
	}
	if req.ShortDescription != "" {
		product.ShortDescription = req.ShortDescription
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.CompareAtPrice > 0 {
		product.CompareAtPrice = req.CompareAtPrice
	}
	if req.CostPerItem > 0 {
		product.CostPerItem = req.CostPerItem
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.LowStockThreshold >= 0 {
		product.LowStockThreshold = req.LowStockThreshold
	}
	if req.TrackInventory != nil {
		product.TrackInventory = *req.TrackInventory
	}
	if req.AllowBackorder != nil {
		product.AllowBackorder = *req.AllowBackorder
	}
	if req.Weight > 0 {
		product.Weight = req.Weight
	}
	if req.Length > 0 {
		product.Length = req.Length
	}
	if req.Width > 0 {
		product.Width = req.Width
	}
	if req.Height > 0 {
		product.Height = req.Height
	}
	if req.MetaTitle != "" {
		product.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		product.MetaDescription = req.MetaDescription
	}
	if req.MetaKeywords != "" {
		product.MetaKeywords = req.MetaKeywords
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}

	if req.CategoryIDs != nil {
		product.Categories = []models.Category{}
		for _, categoryID := range req.CategoryIDs {
			category, err := s.categoryRepo.GetByID(categoryID)
			if err != nil {
				return nil, fmt.Errorf("category with ID %d not found", categoryID)
			}
			product.Categories = append(product.Categories, *category)
		}
	}

	err = s.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productService) DeleteProduct(id uint) error {
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(id)
}

func (s *productService) ListProducts(req *request.ListProductsRequest) (*response.ProductListResponse, error) {
	page := 1
	if req.Page > 0 {
		page = req.Page
	}

	limit := 20
	if req.Limit > 0 {
		limit = req.Limit
	}

	products, total, err := s.productRepo.List(
		page, limit, req.Search, req.VendorID, req.CategoryID,
		req.MinPrice, req.MaxPrice, req.Status, req.IsFeatured, req.IsPublished,
		req.Sort, req.Order,
	)
	if err != nil {
		return nil, err
	}

	vendorIDs := make([]uint, 0, len(products))
	vendorIDMap := make(map[uint]bool)
	for _, product := range products {
		if !vendorIDMap[product.VendorID] {
			vendorIDs = append(vendorIDs, product.VendorID)
			vendorIDMap[product.VendorID] = true
		}
	}

	vendorMap, err := s.vendorService.GetVendorsByIDs(vendorIDs)
	if err != nil {
		vendorMap = make(map[uint]*clients.VendorResponse)
	}

	var productResponses []response.ProductListItemResponse
	for _, product := range products {
		vendorResp, exists := vendorMap[product.VendorID]
		if !exists {
			vendorResp = &clients.VendorResponse{ID: product.VendorID, Name: "Unknown Vendor", Slug: ""}
		}

		item := s.toProductListItemResponse(&product)
		item.Vendor = response.VendorResponse{
			ID:   vendorResp.ID,
			Name: vendorResp.Name,
			Slug: vendorResp.Slug,
		}
		productResponses = append(productResponses, *item)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &response.ProductListResponse{
		Data: productResponses,
		Pagination: response.PaginationResponse{
			Total:      int(total),
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *productService) UpdateProductStatus(id uint, req *request.UpdateProductStatusRequest) (*response.ProductStatusResponse, error) {
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = s.productRepo.UpdateStatus(id, req.Status)
	if err != nil {
		return nil, err
	}

	return &response.ProductStatusResponse{
		Message:   "Product status updated successfully",
		Status:    req.Status,
		UpdatedAt: time.Now(),
	}, nil
}

func (s *productService) PublishProduct(id uint) (*response.ProductPublishResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if product.Status != "active" {
		return nil, errors.New("cannot publish a product that is not active")
	}

	err = s.productRepo.Publish(id)
	if err != nil {
		return nil, err
	}

	return &response.ProductPublishResponse{
		Message:     "Product published successfully",
		IsPublished: true,
		PublishedAt: &time.Time{},
	}, nil
}

func (s *productService) UnpublishProduct(id uint) (*response.ProductPublishResponse, error) {
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = s.productRepo.Unpublish(id)
	if err != nil {
		return nil, err
	}

	return &response.ProductPublishResponse{
		Message:     "Product unpublished successfully",
		IsPublished: false,
		PublishedAt: nil,
	}, nil
}

func (s *productService) toProductResponse(product *models.Product) *response.ProductResponse {
	return &response.ProductResponse{
		ID:                product.ID,
		UUID:              product.UUID.String(),
		VendorID:          product.VendorID,
		Name:              product.Name,
		Slug:              product.Slug,
		SKU:               product.SKU,
		Description:       product.Description,
		ShortDescription:  product.ShortDescription,
		Price:             product.Price,
		CompareAtPrice:    product.CompareAtPrice,
		CostPerItem:       product.CostPerItem,
		Stock:             product.Stock,
		LowStockThreshold: product.LowStockThreshold,
		TrackInventory:    product.TrackInventory,
		AllowBackorder:    product.AllowBackorder,
		Weight:            product.Weight,
		Length:            product.Length,
		Width:             product.Width,
		Height:            product.Height,
		MetaTitle:         product.MetaTitle,
		MetaDescription:   product.MetaDescription,
		MetaKeywords:      product.MetaKeywords,
		Status:            product.Status,
		IsFeatured:        product.IsFeatured,
		IsPublished:       product.IsPublished,
		PublishedAt:       product.PublishedAt,
		ViewCount:         product.ViewCount,
		SoldCount:         product.SoldCount,
		RatingAvg:         product.RatingAvg,
		ReviewCount:       product.ReviewCount,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	}
}

func (s *productService) toProductDetailResponse(product *models.Product) (*response.ProductDetailResponse, error) {
	vendorResp, err := s.vendorService.GetVendorByID(product.VendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vendor details for product %d: %w", product.ID, err)
	}

	vendor := response.VendorResponse{
		ID:   vendorResp.ID,
		Name: vendorResp.Name,
		Slug: vendorResp.Slug,
	}

	var categories []response.CategoryResponse
	for _, category := range product.Categories {
		categories = append(categories, response.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	var images []response.ProductImageResponse
	for _, image := range product.Images {
		images = append(images, response.ProductImageResponse{
			ID:        image.ID,
			ImageURL:  image.ImageURL,
			AltText:   image.AltText,
			SortOrder: image.SortOrder,
			IsPrimary: image.IsPrimary,
		})
	}

	var variants []response.ProductVariantResponse
	for _, variant := range product.Variants {
		variants = append(variants, response.ProductVariantResponse{
			ID:       variant.ID,
			Name:     variant.Name,
			SKU:      variant.SKU,
			Price:    variant.Price,
			Stock:    variant.Stock,
			ImageURL: variant.ImageURL,
			IsActive: variant.IsActive,
		})
	}

	return &response.ProductDetailResponse{
		ID:                product.ID,
		UUID:              product.UUID.String(),
		VendorID:          product.VendorID,
		Vendor:            vendor,
		Name:              product.Name,
		Slug:              product.Slug,
		SKU:               product.SKU,
		Description:       product.Description,
		ShortDescription:  product.ShortDescription,
		Price:             product.Price,
		CompareAtPrice:    product.CompareAtPrice,
		CostPerItem:       product.CostPerItem,
		Stock:             product.Stock,
		LowStockThreshold: product.LowStockThreshold,
		TrackInventory:    product.TrackInventory,
		AllowBackorder:    product.AllowBackorder,
		Weight:            product.Weight,
		Length:            product.Length,
		Width:             product.Width,
		Height:            product.Height,
		MetaTitle:         product.MetaTitle,
		MetaDescription:   product.MetaDescription,
		MetaKeywords:      product.MetaKeywords,
		Status:            product.Status,
		IsFeatured:        product.IsFeatured,
		IsPublished:       product.IsPublished,
		PublishedAt:       product.PublishedAt,
		ViewCount:         product.ViewCount,
		SoldCount:         product.SoldCount,
		RatingAvg:         product.RatingAvg,
		ReviewCount:       product.ReviewCount,
		Categories:        categories,
		Images:            images,
		Variants:          variants,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	}, nil
}

func (s *productService) toProductListItemResponse(product *models.Product) *response.ProductListItemResponse {
	vendor := response.VendorResponse{
		ID:   product.Vendor.ID,
		Name: product.Vendor.Name,
		Slug: product.Vendor.Slug,
	}

	var primaryImage *response.ProductImageResponse
	for _, image := range product.Images {
		if image.IsPrimary {
			primaryImage = &response.ProductImageResponse{
				ID:       image.ID,
				ImageURL: image.ImageURL,
			}
			break
		}
	}

	return &response.ProductListItemResponse{
		ID:             product.ID,
		UUID:           product.UUID.String(),
		VendorID:       product.VendorID,
		Vendor:         vendor,
		Name:           product.Name,
		Slug:           product.Slug,
		SKU:            product.SKU,
		Price:          product.Price,
		CompareAtPrice: product.CompareAtPrice,
		Stock:          product.Stock,
		IsFeatured:     product.IsFeatured,
		IsPublished:    product.IsPublished,
		ViewCount:      product.ViewCount,
		SoldCount:      product.SoldCount,
		RatingAvg:      product.RatingAvg,
		ReviewCount:    product.ReviewCount,
		PrimaryImage:   primaryImage,
		CreatedAt:      product.CreatedAt,
		UpdatedAt:      product.UpdatedAt,
	}
}
