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

type UserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AuthClient interface {
	GetUserByID(id uint) (*UserResponse, error)
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

func (c *authClient) GetUserByID(id uint) (*UserResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("user ID cannot be zero")
	}

	url := fmt.Sprintf("%s/internal/users/%d", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	internalKey := os.Getenv("INTERNAL_SERVICE_KEY")
	if internalKey != "" {
		req.Header.Set("X-Internal-API-Key", internalKey)
	}

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

	var userResp UserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user response: %w", err)
	}

	return &userResp, nil
}
