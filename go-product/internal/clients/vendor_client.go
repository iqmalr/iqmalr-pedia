package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/iqmalr-pedia/go-product/internal/config"
)

type VendorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type VendorClient interface {
	GetVendorByID(id uint) (*VendorResponse, error)
	GetVendorsByIDs(ids []uint) (map[uint]*VendorResponse, error)
}

type vendorClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewVendorClient() VendorClient {
	baseURL := config.AppConfig.VendorServiceURL

	return &vendorClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *vendorClient) GetVendorByID(id uint) (*VendorResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("vendor ID cannot be zero")
	}

	url := fmt.Sprintf("%s/vendors/%d", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	internalKey := os.Getenv("INTERNAL_SERVICE_KEY")
	if internalKey != "" {
		req.Header.Set("X-Internal-Key", internalKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to vendor service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vendor service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var vendorResp VendorResponse
	if err := json.Unmarshal(body, &vendorResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vendor response: %w", err)
	}

	return &vendorResp, nil
}

func (c *vendorClient) GetVendorsByIDs(ids []uint) (map[uint]*VendorResponse, error) {
	if len(ids) == 0 {
		return make(map[uint]*VendorResponse), nil
	}

	idsStr := ""
	for i, id := range ids {
		if i > 0 {
			idsStr += ","
		}
		idsStr += fmt.Sprintf("%d", id)
	}

	url := fmt.Sprintf("%s/vendors?ids=%s", c.baseURL, idsStr)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	internalKey := os.Getenv("INTERNAL_SERVICE_KEY")
	if internalKey != "" {
		req.Header.Set("X-Internal-Key", internalKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to vendor service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vendor service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var vendors []VendorResponse
	if err := json.Unmarshal(body, &vendors); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vendor response: %w", err)
	}

	vendorMap := make(map[uint]*VendorResponse)
	for i := range vendors {
		vendorMap[vendors[i].ID] = &vendors[i]
	}

	return vendorMap, nil
}
