package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/iqmalr-pedia/go-transactions/internal/config"
)

type AddressInfo struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
}

type AuthClient interface {
	GetUserAddressByID(authToken string, addressID uint) (*AddressInfo, error)
}

type authClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient() AuthClient {
	return &authClient{
		baseURL: config.AppConfig.AuthServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *authClient) GetUserAddressByID(authToken string, addressID uint) (*AddressInfo, error) {
	if addressID == 0 {
		return nil, fmt.Errorf("address ID cannot be zero")
	}

	url := fmt.Sprintf("%s/users/me/addresses/%d", c.baseURL, addressID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var address AddressInfo
	if err := json.Unmarshal(body, &address); err != nil {
		return nil, fmt.Errorf("failed to unmarshal address response: %w", err)
	}

	return &address, nil
}
