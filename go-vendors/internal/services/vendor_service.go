package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/iqmalr-pedia/go-vendors/internal/config"
	"github.com/iqmalr-pedia/go-vendors/internal/dto/request"
	"github.com/iqmalr-pedia/go-vendors/internal/dto/response"
	"github.com/iqmalr-pedia/go-vendors/internal/models"
	"github.com/iqmalr-pedia/go-vendors/internal/repositories"
	"github.com/iqmalr-pedia/go-vendors/internal/utils"
)

type VendorService struct {
	vendorRepo repositories.VendorRepositoryInterface
	httpClient utils.HTTPClient
}

func NewVendorService(vendorRepo repositories.VendorRepositoryInterface, client utils.HTTPClient) *VendorService {
	return &VendorService{
		vendorRepo: vendorRepo,
		httpClient: client,
	}
}

func (s *VendorService) validateUserWithAuthService(email string) (*response.UserValidationResponse, error) {
	url := fmt.Sprintf("%s/api/v2/internal/users/by-email?email=%s", config.AppConfig.AuthServiceURL, email)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Internal-API-Key", config.AppConfig.InternalAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("user not found")
	}

	return utils.ParseUserValidationResponse(resp)
}

func (s *VendorService) CreateVendor(userID uint, req *request.CreateVendorRequest) (*response.VendorResponse, error) {
	_, err := s.vendorRepo.FindByOwnerID(userID)
	if err == nil {
		return nil, errors.New("user already has a vendor")
	}

	slug := utils.GenerateSlug(req.Name)
	for {
		_, err := s.vendorRepo.FindBySlug(slug)
		if err != nil {
			break
		}
		// TODO: Implement a better slug uniqueness logic (e.g., add counter)
		slug += "-1"
	}

	vendor := &models.Vendor{
		OwnerID:      userID,
		Name:         req.Name,
		Slug:         slug,
		Description:  req.Description,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		BusinessType: req.BusinessType,
		TaxID:        req.TaxID,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		Status:       "pending",
	}

	if vendor.Country == "" {
		vendor.Country = "Indonesia"
	}

	if err := s.vendorRepo.Create(vendor); err != nil {
		return nil, err
	}

	if err := s.vendorRepo.AddUserToVendor(vendor.ID, userID, "owner"); err != nil {
		// TODO: Handle rollback (delete vendor if adding user fails)
		return nil, err
	}

	return s.mapToVendorResponse(vendor), nil
}

func (s *VendorService) GetVendorByID(id uint) (*response.VendorResponse, error) {
	vendor, err := s.vendorRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("vendor not found")
	}
	return s.mapToVendorResponse(vendor), nil
}

func (s *VendorService) GetVendorBySlug(slug string) (*response.VendorResponse, error) {
	vendor, err := s.vendorRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("vendor not found")
	}
	return s.mapToVendorResponse(vendor), nil
}

func (s *VendorService) UpdateVendor(userID, vendorID uint, userRole string, req *request.UpdateVendorRequest) (*response.VendorResponse, error) {
	vendor, err := s.vendorRepo.FindByID(vendorID)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	isOwner := vendor.OwnerID == userID
	isAdmin := userRole == "admin"

	if !isOwner && !isAdmin {
		return nil, errors.New("unauthorized: you are not the owner or an admin")
	}

	if req.Name != "" {
		vendor.Name = req.Name
	}
	if req.Description != "" {
		vendor.Description = req.Description
	}
	if req.ContactEmail != "" {
		vendor.ContactEmail = req.ContactEmail
	}
	if req.ContactPhone != "" {
		vendor.ContactPhone = req.ContactPhone
	}
	if req.BusinessType != "" {
		vendor.BusinessType = req.BusinessType
	}
	if req.TaxID != "" {
		vendor.TaxID = req.TaxID
	}
	if req.AddressLine1 != "" {
		vendor.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != "" {
		vendor.AddressLine2 = req.AddressLine2
	}
	if req.City != "" {
		vendor.City = req.City
	}
	if req.State != "" {
		vendor.State = req.State
	}
	if req.PostalCode != "" {
		vendor.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		vendor.Country = req.Country
	}

	if err := s.vendorRepo.Update(vendor); err != nil {
		return nil, err
	}

	return s.mapToVendorResponse(vendor), nil
}

func (s *VendorService) UpdateVendorStatus(vendorID uint, status string) (*response.MessageResponse, error) {
	if err := s.vendorRepo.UpdateStatus(vendorID, status); err != nil {
		return nil, err
	}
	return &response.MessageResponse{Message: "Vendor status updated successfully"}, nil
}

func (s *VendorService) ListVendors(req *request.ListVendorsRequest) (*response.VendorListResponse, error) {
	vendors, total, err := s.vendorRepo.FindWithPagination(req.Page, req.Limit, req.Search, req.Status, req.Sort, req.Order)
	if err != nil {
		return nil, err
	}

	vendorListItems := make([]response.VendorListItem, len(vendors))
	for i, vendor := range vendors {
		vendorListItems[i] = response.VendorListItem{
			ID:            vendor.ID,
			UUID:          vendor.UUID.String(),
			Name:          vendor.Name,
			Slug:          vendor.Slug,
			LogoUrl:       vendor.LogoUrl,
			ContactEmail:  vendor.ContactEmail,
			Status:        vendor.Status,
			TotalProducts: vendor.TotalProducts,
			RatingAvg:     vendor.RatingAvg,
			TotalReviews:  vendor.TotalReviews,
			CreatedAt:     vendor.CreatedAt,
		}
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &response.VendorListResponse{
		Data: vendorListItems,
		Pagination: response.PaginationResponse{
			TotalPages: totalPages,
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
		},
	}, nil
}

func (s *VendorService) mapToVendorResponse(vendor *models.Vendor) *response.VendorResponse {
	return &response.VendorResponse{
		ID:            vendor.ID,
		UUID:          vendor.UUID.String(),
		OwnerID:       vendor.OwnerID,
		Name:          vendor.Name,
		Slug:          vendor.Slug,
		Description:   vendor.Description,
		LogoUrl:       vendor.LogoUrl,
		BannerUrl:     vendor.BannerUrl,
		ContactEmail:  vendor.ContactEmail,
		ContactPhone:  vendor.ContactPhone,
		Status:        vendor.Status,
		BusinessType:  vendor.BusinessType,
		TaxID:         vendor.TaxID,
		AddressLine1:  vendor.AddressLine1,
		AddressLine2:  vendor.AddressLine2,
		City:          vendor.City,
		State:         vendor.State,
		PostalCode:    vendor.PostalCode,
		Country:       vendor.Country,
		TotalProducts: vendor.TotalProducts,
		TotalSales:    vendor.TotalSales,
		RatingAvg:     vendor.RatingAvg,
		TotalReviews:  vendor.TotalReviews,
		VerifiedAt:    vendor.VerifiedAt,
		ApprovedAt:    vendor.ApprovedAt,
		CreatedAt:     vendor.CreatedAt,
		UpdatedAt:     vendor.UpdatedAt,
	}
}
