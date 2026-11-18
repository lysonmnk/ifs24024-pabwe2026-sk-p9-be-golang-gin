package post

import (
	"math"
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)

func GenerateSlug(title string) string {
	lower := strings.ToLower(title)
	slug := nonAlphanumericRegex.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "artikel"
	}
	return slug
}

func CalculateReadingTime(content string) int {
	words := strings.Fields(content)
	wordCount := len(words)
	if wordCount == 0 {
		return 1
	}

	return int(math.Ceil(float64(wordCount) / 200.0))
}
