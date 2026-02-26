package utils

import (
	"regexp"
	"strings"
	"testing"
)

func TestGenerateSKU_Format(t *testing.T) {
	sku := GenerateSKU(1, "Kaos Polos Premium")
	pattern := regexp.MustCompile(`^V001-KAO-\d{6}-\d{6}$`)
	if !pattern.MatchString(sku) {
		t.Errorf("SKU format mismatch, got %q", sku)
	}
}

func TestGenerateSKU_VendorPrefix(t *testing.T) {
	tests := []struct {
		vendorID uint
		expected string
	}{
		{1, "V001-"},
		{42, "V042-"},
		{999, "V999-"},
	}
	for _, tt := range tests {
		sku := GenerateSKU(tt.vendorID, "Test Product")
		if !strings.HasPrefix(sku, tt.expected) {
			t.Errorf("vendorID=%d: expected prefix %q, got %q", tt.vendorID, tt.expected, sku)
		}
	}
}

func TestGenerateSKU_ShortName(t *testing.T) {
	sku := GenerateSKU(1, "AB")
	parts := strings.Split(sku, "-")
	if parts[1] != "ABX" {
		t.Errorf("short name should be padded with X, got product prefix %q", parts[1])
	}
}

func TestGenerateSKU_EmptyName(t *testing.T) {
	sku := GenerateSKU(1, "")
	parts := strings.Split(sku, "-")
	if parts[1] != "XXX" {
		t.Errorf("empty name should produce XXX prefix, got %q", parts[1])
	}
}

func TestGenerateSKU_SpecialCharsInName(t *testing.T) {
	sku := GenerateSKU(1, "!@# Hello World")
	parts := strings.Split(sku, "-")
	if parts[1] != "HEL" {
		t.Errorf("special chars should be skipped, expected HEL, got %q", parts[1])
	}
}

func TestGenerateSKU_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		sku := GenerateSKU(1, "Test")
		if seen[sku] {
			t.Errorf("duplicate SKU generated: %q", sku)
		}
		seen[sku] = true
	}
}

func TestExtractProductPrefix_Unicode(t *testing.T) {
	prefix := extractProductPrefix("日本語テスト", 3)
	if len(prefix) < 3 {
		t.Errorf("prefix should be at least 3 chars, got %q", prefix)
	}
}

func TestCryptoRandom_Length(t *testing.T) {
	for _, digits := range []int{4, 6, 8} {
		result := cryptoRandom(digits)
		if len(result) != digits {
			t.Errorf("expected length %d, got %d (%q)", digits, len(result), result)
		}
	}
}
