package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/iqmalr-pedia/go-transactions/internal/config"
)

type ProductInfo struct {
	ID             uint        `json:"id"`
	VendorID       uint        `json:"vendor_id"`
	Name           string      `json:"name"`
	Slug           string      `json:"slug"`
	SKU            string      `json:"sku"`
	Price          float64     `json:"price"`
	Stock          int         `json:"stock"`
	TrackInventory bool        `json:"track_inventory"`
	AllowBackorder bool        `json:"allow_backorder"`
	IsPublished    bool        `json:"is_published"`
	Vendor         *VendorInfo `json:"vendor,omitempty"`
	PrimaryImage   *Image      `json:"primary_image,omitempty"`
}

type VendorInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Image struct {
	ID       uint   `json:"id"`
	ImageURL string `json:"image_url"`
}

type VariantInfo struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	IsActive bool    `json:"is_active"`
}

type ProductClient interface {
	GetProductByID(id uint) (*ProductInfo, error)
	GetVariantByID(productID, variantID uint) (*VariantInfo, error)
}

type productClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductClient() ProductClient {
	baseURL := config.AppConfig.ProductServiceURL

	return &productClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *productClient) GetProductByID(id uint) (*ProductInfo, error) {
	if id == 0 {
		return nil, fmt.Errorf("product ID cannot be zero")
	}

	url := fmt.Sprintf("%s/products/%d", c.baseURL, id)

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
		return nil, fmt.Errorf("failed to execute request to product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("product service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var product ProductInfo
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product response: %w", err)
	}

	return &product, nil
}

func (c *productClient) GetVariantByID(productID, variantID uint) (*VariantInfo, error) {
	if productID == 0 || variantID == 0 {
		return nil, fmt.Errorf("product ID and variant ID cannot be zero")
	}

	url := fmt.Sprintf("%s/products/%d/variants/%d", c.baseURL, productID, variantID)

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
		return nil, fmt.Errorf("failed to execute request to product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("product service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var variant VariantInfo
	if err := json.Unmarshal(body, &variant); err != nil {
		return nil, fmt.Errorf("failed to unmarshal variant response: %w", err)
	}

	return &variant, nil
}
