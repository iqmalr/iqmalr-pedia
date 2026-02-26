package utils

import (
	"testing"
)

func TestGenerateProductSKU(t *testing.T) {
	tests := []struct {
		name         string
		vendorID     uint
		categoryCode string
		sequence     uint
		want         string
	}{
		{
			name:         "Standard case",
			vendorID:     1,
			categoryCode: "T-Shirt",
			sequence:     42,
			want:         "V001-TSH-00042",
		},
		{
			name:         "Vendor with leading zeros",
			vendorID:     15,
			categoryCode: "Electronics",
			sequence:     1,
			want:         "V015-ELE-00001",
		},
		{
			name:         "Special characters in category",
			vendorID:     99,
			categoryCode: "Bottoms & Jeans!",
			sequence:     999,
			want:         "V099-BOT-00999",
		},
		{
			name:         "Empty category code",
			vendorID:     5,
			categoryCode: "",
			sequence:     100,
			want:         "V005-XXX-00100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateProductSKU(tt.vendorID, tt.categoryCode, tt.sequence); got != tt.want {
				t.Errorf("GenerateProductSKU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateVariantSKU(t *testing.T) {
	tests := []struct {
		name        string
		parentSKU   string
		variantName string
		want        string
	}{
		{
			name:        "Simple variant",
			parentSKU:   "V001-TSH-00042",
			variantName: "White XL",
			want:        "V001-TSH-00042-WHITXL",
		},
		{
			name:        "Variant with special chars",
			parentSKU:   "V015-ELE-00001",
			variantName: "Color: Red / Size: M",
			want:        "V015-ELE-00001-COLOR",
		},
		{
			name:        "Short variant name",
			parentSKU:   "V099-BOT-00999",
			variantName: "S",
			want:        "V099-BOT-00999-SXX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateVariantSKU(tt.parentSKU, tt.variantName); got != tt.want {
				t.Errorf("GenerateVariantSKU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeCode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		want      string
	}{
		{"Normal string", "T-Shirt", 3, "TSH"},
		{"With spaces", "Long John", 3, "LON"},
		{"With special chars", "Bottoms & Jeans!", 3, "BOT"},
		{"Empty string", "", 3, "XXX"},
		{"Short string", "A", 3, "AXX"},
		{"Long string", "VeryLongCategoryName", 4, "VERY"},
		{"With numbers", "Product 2024", 5, "PRODU"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeCode(tt.input, tt.maxLength); got != tt.want {
				t.Errorf("sanitizeCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
