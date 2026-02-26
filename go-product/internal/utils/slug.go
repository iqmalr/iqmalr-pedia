package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const maxSlugBase = 60

func GenerateSlug(name string, uniqueID string) string {
	slug := normalizeUnicode(name)

	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")

	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	slug = strings.Trim(slug, "-")

	if len(slug) > maxSlugBase {
		slug = slug[:maxSlugBase]
		slug = strings.TrimRight(slug, "-")
	}

	if uniqueID != "" {
		uniqueID = strings.ToLower(uniqueID)
		if len(uniqueID) > 8 {
			uniqueID = uniqueID[:8]
		}

		if slug == "" {
			slug = uniqueID
		} else {
			slug = slug + "-" + uniqueID
		}
	}

	return slug
}

func normalizeUnicode(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return result
}
