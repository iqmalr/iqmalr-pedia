package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/models"
	"github.com/iqmalr-pedia/go-product/internal/repositories"
)

type VariantService interface {
	CreateVariant(productID uint, req *request.CreateVariantRequest) (*response.VariantDetailResponse, error)
	GetVariantsByProductID(productID uint) (*response.VariantListResponse, error)
	UpdateVariant(productID, variantID uint, req *request.UpdateVariantRequest) (*response.VariantDetailResponse, error)
	DeleteVariant(productID, variantID uint) error
}

type variantService struct {
	variantRepo repositories.VariantRepository
	productRepo repositories.ProductRepository
}

func NewVariantService(variantRepo repositories.VariantRepository, productRepo repositories.ProductRepository) VariantService {
	return &variantService{
		variantRepo: variantRepo,
		productRepo: productRepo,
	}
}

func (s *variantService) CreateVariant(productID uint, req *request.CreateVariantRequest) (*response.VariantDetailResponse, error) {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	sku := req.SKU
	if sku == "" {
		sku = generateVariantSKU(product.VendorID, req.Name)
	} else {
		exists, err := s.variantRepo.SKUExistsForProduct(productID, sku, 0)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("SKU '%s' already exists for this product", sku)
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	variant := &models.ProductVariant{
		ProductID: productID,
		Name:      req.Name,
		SKU:       sku,
		Price:     req.Price,
		Stock:     req.Stock,
		ImageURL:  req.ImageURL,
		IsActive:  isActive,
	}

	if err := s.variantRepo.Create(variant); err != nil {
		return nil, err
	}

	return s.toVariantDetailResponse(variant), nil
}

func (s *variantService) GetVariantsByProductID(productID uint) (*response.VariantListResponse, error) {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	variants, err := s.variantRepo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}

	var items []response.VariantDetailResponse
	for _, v := range variants {
		items = append(items, *s.toVariantDetailResponse(&v))
	}

	return &response.VariantListResponse{Data: items}, nil
}

func (s *variantService) UpdateVariant(productID, variantID uint, req *request.UpdateVariantRequest) (*response.VariantDetailResponse, error) {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	variant, err := s.variantRepo.GetByID(variantID)
	if err != nil {
		return nil, err
	}

	if variant.ProductID != productID {
		return nil, fmt.Errorf("variant does not belong to this product")
	}

	if req.Name != "" {
		variant.Name = req.Name
	}

	if req.SKU != "" && req.SKU != variant.SKU {
		exists, err := s.variantRepo.SKUExistsForProduct(productID, req.SKU, variantID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("SKU '%s' already exists for this product", req.SKU)
		}
		variant.SKU = req.SKU
	}

	if req.Price > 0 {
		variant.Price = req.Price
	}

	if req.Stock != nil {
		variant.Stock = *req.Stock
	}

	if req.ImageURL != "" {
		variant.ImageURL = req.ImageURL
	}

	if req.IsActive != nil {
		variant.IsActive = *req.IsActive
	}

	if err := s.variantRepo.Update(variant); err != nil {
		return nil, err
	}

	return s.toVariantDetailResponse(variant), nil
}

func (s *variantService) DeleteVariant(productID, variantID uint) error {
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	variant, err := s.variantRepo.GetByID(variantID)
	if err != nil {
		return err
	}

	if variant.ProductID != productID {
		return fmt.Errorf("variant does not belong to this product")
	}

	return s.variantRepo.Delete(variantID)
}

func (s *variantService) toVariantDetailResponse(v *models.ProductVariant) *response.VariantDetailResponse {
	return &response.VariantDetailResponse{
		ID:        v.ID,
		ProductID: v.ProductID,
		Name:      v.Name,
		SKU:       v.SKU,
		Price:     v.Price,
		Stock:     v.Stock,
		ImageURL:  v.ImageURL,
		IsActive:  v.IsActive,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func generateVariantSKU(vendorID uint, variantName string) string {
	now := time.Now()
	vendorPrefix := fmt.Sprintf("V%03d", vendorID)

	cleaned := strings.ReplaceAll(variantName, " ", "")
	namePrefix := ""
	if len(cleaned) >= 3 {
		namePrefix = strings.ToUpper(cleaned[:3])
	} else {
		namePrefix = strings.ToUpper(cleaned)
		for len(namePrefix) < 3 {
			namePrefix += "X"
		}
	}

	dateComponent := now.Format("060102")
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(10000))
	randomComponent := fmt.Sprintf("%04d", randomNum.Int64())

	return fmt.Sprintf("%s-VAR-%s-%s-%s", vendorPrefix, namePrefix, dateComponent, randomComponent)
}
