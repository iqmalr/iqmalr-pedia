package services

import (
	"errors"

	"github.com/iqmalr-pedia/go-vendors/internal/dto/request"
	"github.com/iqmalr-pedia/go-vendors/internal/dto/response"
	"github.com/iqmalr-pedia/go-vendors/internal/models"
	"github.com/iqmalr-pedia/go-vendors/internal/repositories"
	"gorm.io/gorm"
)

type VendorAccountBankService struct {
	vendorRepo            repositories.VendorRepositoryInterface
	vendorAccountBankRepo repositories.VendorAccountBankRepositoryInterface
}

func NewVendorAccountBankService(v repositories.VendorRepositoryInterface, va repositories.VendorAccountBankRepositoryInterface) *VendorAccountBankService {
	return &VendorAccountBankService{
		vendorRepo:            v,
		vendorAccountBankRepo: va,
	}
}

func (v *VendorAccountBankService) mapAccountBankResponse(vendor *models.VendorBankAccount) *response.VendorAccountBankResponse {
	return &response.VendorAccountBankResponse{
		ID:            vendor.ID,
		VendorID:      vendor.VendorID,
		BankName:      vendor.BankName,
		AccountNumber: vendor.AccountNumber,
		AccountHolder: vendor.AccountHolder,
		IsVerified:    vendor.IsVerified,
		IsPrimary:     vendor.IsPrimary,
		CreatedAt:     vendor.CreatedAt,
		UpdatedAt:     vendor.UpdatedAt,
	}
}
func (v *VendorAccountBankService) mapArrayAccountBankResponses(vendors *[]models.VendorBankAccount) *[]response.VendorAccountBankResponse {

	responses := make([]response.VendorAccountBankResponse, 0)

	for _, vendor := range *vendors {
		responses = append(responses, response.VendorAccountBankResponse{
			ID:            vendor.ID,
			VendorID:      vendor.VendorID,
			BankName:      vendor.BankName,
			AccountNumber: vendor.AccountNumber,
			AccountHolder: vendor.AccountHolder,
			IsVerified:    vendor.IsVerified,
			IsPrimary:     vendor.IsPrimary,
			CreatedAt:     vendor.CreatedAt,
			UpdatedAt:     vendor.UpdatedAt,
		})
	}

	return &responses
}

func (v *VendorAccountBankService) CreateAccountBank(VendorID uint, req *request.CreateAccountBank) (*response.VendorAccountBankResponse, error) {
	_, err := v.vendorRepo.FindByID(VendorID)
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("Vendor is not found")
	}

	vendorAccountBank := &models.VendorBankAccount{
		VendorID:      VendorID,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		AccountHolder: req.AccountHolder,
		IsPrimary:     req.IsPrimary,
	}

	if err := v.vendorAccountBankRepo.CreateAccountBank(vendorAccountBank); err != nil {

	}
	return v.mapAccountBankResponse(vendorAccountBank), nil
}

func (v *VendorAccountBankService) GetAccountBankByID(VendorID uint) (*[]response.VendorAccountBankResponse, error) {
	_, err := v.vendorRepo.FindByID(VendorID)
	if err != nil {
		return nil, errors.New("Bank account not found")
	}

	value, err := v.vendorAccountBankRepo.GetAccountBankByVendorID(VendorID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return v.mapArrayAccountBankResponses(value), nil
}

func (v *VendorAccountBankService) UpdateAccountBank(AccountID uint, req *request.UpdateAccountBank) (*response.VendorAccountBankResponse, error) {
	value, err := v.vendorAccountBankRepo.GetAccountBankByID(AccountID)
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("This account not found")
	}

	vendorAccountBank := &models.VendorBankAccount{
		ID:            value.ID,
		VendorID:      value.VendorID,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		AccountHolder: req.AccountHolder,
		IsPrimary:     req.IsPrimary,
	}

	if err := v.vendorAccountBankRepo.UpdateAccountBank(AccountID, vendorAccountBank); err != nil {

	}
	return v.mapAccountBankResponse(vendorAccountBank), nil
}

func (v *VendorAccountBankService) DeleteAccountBank(VendorID, AccountID uint) (*response.MessageResponse, error) {
	_, err := v.vendorAccountBankRepo.GetAccountBankByID(AccountID)
	if err != nil {
		return nil, errors.New("Bank account not found")
	}

	if err := v.vendorAccountBankRepo.DeleteAccountBank(VendorID, AccountID); err != nil {
		return nil, err
	}

	return &response.MessageResponse{Message: "Bank account deleted successfully"}, nil
}
