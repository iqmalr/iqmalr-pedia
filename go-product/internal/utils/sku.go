package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
	"unicode"
)

// GenerateSKU creates a SKU with format: V{vendorID}-{productPrefix}-{YYMMDD}-{random6digit}
// Example: V001-KAO-260226-482917
func GenerateSKU(vendorID uint, productName string) string {
	vendorPrefix := fmt.Sprintf("V%03d", vendorID)
	productPrefix := extractProductPrefix(productName, 3)
	dateComponent := time.Now().Format("060102")
	randomComponent := cryptoRandom(6)

	return fmt.Sprintf("%s-%s-%s-%s", vendorPrefix, productPrefix, dateComponent, randomComponent)
}

func extractProductPrefix(name string, length int) string {
	var letters []rune
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			letters = append(letters, unicode.ToUpper(r))
		}
		if len(letters) >= length {
			break
		}
	}

	prefix := string(letters)
	for len(prefix) < length {
		prefix += "X"
	}

	return prefix
}

func cryptoRandom(digits int) string {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return strings.Repeat("0", digits)
	}
	return fmt.Sprintf("%0*d", digits, n.Int64())
}
