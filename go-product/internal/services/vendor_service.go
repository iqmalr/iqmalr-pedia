package services

import (
	"github.com/iqmalr-pedia/go-product/internal/clients"
)

type VendorService interface {
	GetVendorByID(id uint) (*clients.VendorResponse, error)
	GetVendorsByIDs(ids []uint) (map[uint]*clients.VendorResponse, error)
}

type vendorService struct {
	client clients.VendorClient
}

func NewVendorService(client clients.VendorClient) VendorService {
	return &vendorService{
		client: client,
	}
}

func (s *vendorService) GetVendorByID(id uint) (*clients.VendorResponse, error) {
	return s.client.GetVendorByID(id)
}

func (s *vendorService) GetVendorsByIDs(ids []uint) (map[uint]*clients.VendorResponse, error) {
	return s.client.GetVendorsByIDs(ids)
}
