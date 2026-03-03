package utils

import "testing"

func TestGenerateSlug_BasicFormat(t *testing.T) {
	slug := GenerateSlug("Kaos Polos Premium", "a1b2c3d4e5")
	expected := "kaos-polos-premium-a1b2c3d4"
	if slug != expected {
		t.Errorf("expected %q, got %q", expected, slug)
	}
}

func TestGenerateSlug_SameNameDifferentID(t *testing.T) {
	slug1 := GenerateSlug("Kaos Polos", "aaaaaaaa")
	slug2 := GenerateSlug("Kaos Polos", "bbbbbbbb")
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

func TestGenerateSlug_NoLeadingDash(t *testing.T) {
	slug := GenerateSlug("---Hello", "abc12345")
	if slug[0] == '-' {
		t.Errorf("slug has leading dash: %q", slug)
	}
}

func TestGenerateSlug_SpecialCharactersRemoved(t *testing.T) {
	slug := GenerateSlug("Product @#$% Name!", "abcd1234")
	expected := "product-name-abcd1234"
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

func TestGenerateSlug_Lowercase(t *testing.T) {
	slug := GenerateSlug("UPPERCASE NAME", "ABC12345")
	expected := "uppercase-name-abc12345"
	if slug != expected {
		t.Errorf("expected %q, got %q", expected, slug)
	}
}

func TestGenerateSlug_EmptyUniqueID(t *testing.T) {
	slug := GenerateSlug("Test Product", "")
	expected := "test-product"
	if slug != expected {
		t.Errorf("expected %q, got %q", expected, slug)
	}
}

func TestGenerateSlug_NumbersPreserved(t *testing.T) {
	slug := GenerateSlug("Item 123 Pack", "aabb1122")
	expected := "item-123-pack-aabb1122"
	if slug != expected {
		t.Errorf("expected %q, got %q", expected, slug)
	}
}
