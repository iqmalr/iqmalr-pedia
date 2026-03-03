package services

import (
	"errors"

	"github.com/iqmalr-pedia/go-auth/v2/internal/dto/request"
	"github.com/iqmalr-pedia/go-auth/v2/internal/dto/response"
	"github.com/iqmalr-pedia/go-auth/v2/internal/models"
	"github.com/iqmalr-pedia/go-auth/v2/internal/repositories"
)

type AddressService struct {
	addressRepo *repositories.AddressRepository
}

func NewAddressService(addressRepo *repositories.AddressRepository) *AddressService {
	return &AddressService{addressRepo: addressRepo}
}

func (s *AddressService) GetAddresses(userID uint) (*response.AddressListResponse, error) {
	addresses, err := s.addressRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	addressResponses := make([]response.AddressResponse, len(addresses))
	for i, addr := range addresses {
		addressResponses[i] = response.AddressResponse{
			ID:            addr.ID,
			Label:         addr.Label,
			RecipientName: addr.RecipientName,
			Phone:         addr.Phone,
			AddressLine1:  addr.AddressLine1,
			AddressLine2:  addr.AddressLine2,
			City:          addr.City,
			State:         addr.State,
			PostalCode:    addr.PostalCode,
			Country:       addr.Country,
			IsDefault:     addr.IsDefault,
			Latitude:      addr.Latitude,
			Longitude:     addr.Longitude,
			CreatedAt:     addr.CreatedAt,
			UpdatedAt:     addr.UpdatedAt,
		}
	}

	return &response.AddressListResponse{Data: addressResponses}, nil
}

func (s *AddressService) CreateAddress(userID uint, req *request.CreateAddressRequest) (*response.AddressResponse, error) {
	address := &models.UserAddress{
		UserID:        userID,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		City:          req.City,
		State:         req.State,
		PostalCode:    req.PostalCode,
		Country:       req.Country,
		IsDefault:     req.IsDefault != nil && *req.IsDefault,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
	}

	if address.Country == "" {
		address.Country = "Indonesia"
	}

	if err := s.addressRepo.Create(address); err != nil {
		return nil, err
	}

	return &response.AddressResponse{
		ID:            address.ID,
		Label:         address.Label,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		AddressLine1:  address.AddressLine1,
		AddressLine2:  address.AddressLine2,
		City:          address.City,
		State:         address.State,
		PostalCode:    address.PostalCode,
		Country:       address.Country,
		IsDefault:     address.IsDefault,
		Latitude:      address.Latitude,
		Longitude:     address.Longitude,
		CreatedAt:     address.CreatedAt,
		UpdatedAt:     address.UpdatedAt,
	}, nil
}

func (s *AddressService) GetAddressByID(userID, addressID uint) (*response.AddressResponse, error) {
	address, err := s.addressRepo.FindByIDAndUserID(addressID, userID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("address not found")
	}

	return &response.AddressResponse{
		ID:            address.ID,
		Label:         address.Label,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		AddressLine1:  address.AddressLine1,
		AddressLine2:  address.AddressLine2,
		City:          address.City,
		State:         address.State,
		PostalCode:    address.PostalCode,
		Country:       address.Country,
		IsDefault:     address.IsDefault,
		Latitude:      address.Latitude,
		Longitude:     address.Longitude,
		CreatedAt:     address.CreatedAt,
		UpdatedAt:     address.UpdatedAt,
	}, nil
}

func (s *AddressService) UpdateAddress(userID, addressID uint, req *request.UpdateAddressRequest) (*response.AddressResponse, error) {
	address, err := s.addressRepo.FindByIDAndUserID(addressID, userID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("address not found")
	}

	if req.Label != "" {
		address.Label = req.Label
	}
	if req.RecipientName != "" {
		address.RecipientName = req.RecipientName
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.AddressLine1 != "" {
		address.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != "" {
		address.AddressLine2 = req.AddressLine2
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.State != "" {
		address.State = req.State
	}
	if req.PostalCode != "" {
		address.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		address.Country = req.Country
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}
	if req.Latitude != nil {
		address.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		address.Longitude = req.Longitude
	}

	if err := s.addressRepo.Update(address); err != nil {
		return nil, err
	}

	return s.GetAddressByID(userID, addressID)
}

func (s *AddressService) DeleteAddress(userID, addressID uint) error {
	address, err := s.addressRepo.FindByIDAndUserID(addressID, userID)
	if err != nil {
		return err
	}
	if address == nil {
		return errors.New("address not found")
	}

	return s.addressRepo.Delete(addressID, userID)
}

func (s *AddressService) SetDefaultAddress(userID, addressID uint) error {
	address, err := s.addressRepo.FindByIDAndUserID(addressID, userID)
	if err != nil {
		return err
	}
	if address == nil {
		return errors.New("address not found")
	}

	return s.addressRepo.SetDefault(addressID, userID)
}
