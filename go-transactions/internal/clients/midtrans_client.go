package clients

import (
	"encoding/json"
	"fmt"

	"github.com/iqmalr-pedia/go-transactions/internal/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransClient interface {
	ChargeBankTransfer(orderID string, grossAmount int64, bank string, customer *midtrans.CustomerDetails) (*coreapi.ChargeResponse, error)
	ChargeGopay(orderID string, grossAmount int64, customer *midtrans.CustomerDetails) (*coreapi.ChargeResponse, error)
	ChargeShopeePay(orderID string, grossAmount int64, customer *midtrans.CustomerDetails) (*coreapi.ChargeResponse, error)
	CheckTransaction(orderID string) (*coreapi.TransactionStatusResponse, error)
}

type midtransClient struct {
	client coreapi.Client
}

func NewMidtransClient() MidtransClient {
	env := midtrans.Sandbox
	if config.AppConfig.MidtransEnv == "production" {
		env = midtrans.Production
	}
	var c coreapi.Client
	c.New(config.AppConfig.MidtransServerKey, env)
	return &midtransClient{client: c}
}

func (m *midtransClient) charge(req *coreapi.ChargeReq) (*coreapi.ChargeResponse, error) {
	res, err := m.client.ChargeTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("midtrans charge: %w", err)
	}
	if res != nil && res.StatusCode != "201" {
		return nil, fmt.Errorf("midtrans: %s - %s", res.StatusCode, res.StatusMessage)
	}
	return res, nil
}

func (m *midtransClient) ChargeBankTransfer(orderID string, grossAmount int64, bank string, customer *midtrans.CustomerDetails) (*coreapi.ChargeResponse, error) {
	bankVal := midtrans.Bank(bank)
	if bankVal == "" {
		bankVal = midtrans.BankBca
	}
	req := &coreapi.ChargeReq{
		PaymentType:        coreapi.PaymentTypeBankTransfer,
		TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: grossAmount},
		CustomerDetails:    customer,
		BankTransfer:       &coreapi.BankTransferDetails{Bank: bankVal},
	}
	return m.charge(req)
}

func (m *midtransClient) ChargeGopay(orderID string, grossAmount int64, customer *midtrans.CustomerDetails) (*coreapi.ChargeResponse, error) {
	req := &coreapi.ChargeReq{
		PaymentType:        coreapi.PaymentTypeGopay,
		TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: grossAmount},
		CustomerDetails:    customer,
		Gopay:              &coreapi.GopayDetails{},
	}
	return m.charge(req)
}

func (m *midtransClient) ChargeShopeePay(orderID string, grossAmount int64, customer *midtrans.CustomerDetails) (*coreapi.ChargeResponse, error) {
	req := &coreapi.ChargeReq{
		PaymentType:        coreapi.PaymentTypeShopeepay,
		TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: grossAmount},
		CustomerDetails:    customer,
		ShopeePay:          &coreapi.ShopeePayDetails{},
	}
	return m.charge(req)
}

func (m *midtransClient) CheckTransaction(orderID string) (*coreapi.TransactionStatusResponse, error) {
	res, err := m.client.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf("midtrans check: %w", err)
	}
	return res, nil
}

func BuildPaymentInstructions(res *coreapi.ChargeResponse) (string, error) {
	if res == nil {
		return "{}", nil
	}
	instructions := map[string]interface{}{
		"transaction_id":   res.TransactionID,
		"expiry_time":      res.ExpiryTime,
		"payment_type":     res.PaymentType,
		"transaction_time": res.TransactionTime,
	}
	if len(res.VaNumbers) > 0 {
		instructions["va_numbers"] = res.VaNumbers
		instructions["bank"] = res.Bank
	}
	if res.PermataVaNumber != "" {
		instructions["permata_va_number"] = res.PermataVaNumber
		instructions["bank"] = "permata"
	}
	if len(res.Actions) > 0 {
		instructions["actions"] = res.Actions
	}
	if res.PaymentCode != "" {
		instructions["payment_code"] = res.PaymentCode
	}
	if res.Store != "" {
		instructions["store"] = res.Store
	}
	data, err := json.Marshal(instructions)
	if err != nil {
		return "{}", err
	}
	return string(data), nil
}
