package utils

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
type CustomHTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *CustomHTTPClient {
	return &CustomHTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CustomHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func CopyHeaders(src, dst http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func CopyResponse(w http.ResponseWriter, resp *http.Response) error {
	CopyHeaders(resp.Header, w.Header())

	w.WriteHeader(resp.StatusCode)

	_, err := io.Copy(w, resp.Body)
	return err
}

func CreateProxyRequest(c *gin.Context, targetURL string) (*http.Request, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	CopyHeaders(c.Request.Header, req.Header)

	req.Header.Set("X-Forwarded-For", c.ClientIP())
	req.Header.Set("X-Forwarded-Host", c.Request.Host)
	req.Header.Set("X-Forwarded-Proto", c.Request.Proto)

	removeHopByHopHeaders(req.Header)

	return req, nil
}

func removeHopByHopHeaders(header http.Header) {
	hopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	for _, h := range hopHeaders {
		header.Del(h)
	}
}
