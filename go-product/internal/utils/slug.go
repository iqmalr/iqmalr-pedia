package utils

import "strings"

func GenerateSlug(name string, uniqueID string) string {
	slug := strings.ToLower(name)
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

	if uniqueID != "" {
		uniqueID = strings.ToLower(uniqueID)
		if len(uniqueID) > 8 {
			uniqueID = uniqueID[:8]
		}
		slug = slug + "-" + uniqueID
	}

	return slug
}
