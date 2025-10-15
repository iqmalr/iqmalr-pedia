// Tambahkan di go-auth/v2/internal/utils/token.go
package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Handle error appropriately in production
	}
	return base64.URLEncoding.EncodeToString(b)
}
