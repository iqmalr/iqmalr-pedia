package utils

import "testing"

func TestGenerateSlug_BasicFormat(t *testing.T) {
	slug := GenerateSlug("Toko Baju Online", "a1b2c3d4e5")
	expected := "toko-baju-online-a1b2c3d4"
	if slug != expected {
		t.Errorf("expected %q, got %q", expected, slug)
	}
}

func TestGenerateSlug_SameNameDifferentID(t *testing.T) {
	slug1 := GenerateSlug("My Vendor", "aaaaaaaa")
	slug2 := GenerateSlug("My Vendor", "bbbbbbbb")
	if slug1 == slug2 {
		t.Error("same name with different UUID should produce different slugs")
	}
}

func TestGenerateSlug_NoDoubleDash(t *testing.T) {
	slug := GenerateSlug("Hello   World!!!", "abc12345")
	for i := 0; i < len(slug)-1; i++ {
		if slug[i] == '-' && slug[i+1] == '-' {
			t.Errorf("slug contains double dash: %q", slug)
			break
		}
	}
}

func TestGenerateSlug_NoTrailingDash(t *testing.T) {
	slug := GenerateSlug("Hello World---", "abc12345")
	if slug[len(slug)-1] == '-' {
		t.Errorf("slug has trailing dash: %q", slug)
	}
}

func TestGenerateSlug_SpecialCharactersRemoved(t *testing.T) {
	slug := GenerateSlug("Vendor @#$% Name!", "abcd1234")
	expected := "vendor-name-abcd1234"
	if slug != expected {
		t.Errorf("expected %q, got %q", expected, slug)
	}
}

func TestGenerateSlug_UniqueIDTruncatedTo8(t *testing.T) {
	slug := GenerateSlug("Test", "abcdefghijklmnop")
	expected := "test-abcdefgh"
	if slug != expected {
		t.Errorf("expected %q, got %q", expected, slug)
	}
}

func TestGenerateSlug_EmptyUniqueID(t *testing.T) {
	slug := GenerateSlug("Test Vendor", "")
	expected := "test-vendor"
	if slug != expected {
		t.Errorf("expected %q, got %q", expected, slug)
	}
}
