package utils

import (
	"fmt"
	"strings"
	"unicode"
)

func GenerateProductSKU(vendorID uint, categoryCode string, sequence uint) string {
	vendorCode := fmt.Sprintf("V%03d", vendorID)
	catCode := sanitizeCode(categoryCode, 3)
	seqCode := fmt.Sprintf("%05d", sequence)

	return fmt.Sprintf("%s-%s-%s", vendorCode, catCode, seqCode)
}

func GenerateVariantSKU(parentSKU string, variantName string) string {
	variantCode := sanitizeCode(variantName, 6)
	return fmt.Sprintf("%s-%s", parentSKU, variantCode)
}

func sanitizeCode(s string, maxLength int) string {
	var result []rune
	for _, r := range strings.ToUpper(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, r)
		}
		if len(result) >= maxLength {
			break
		}
	}

	for len(result) < 3 {
		result = append(result, 'X')
	}

	return string(result)
}
