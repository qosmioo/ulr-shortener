package shortener

import (
	"testing"
)

func TestGenerateShortURL(t *testing.T) {
	shortURL := GenerateShortURL()

	if len(shortURL) != 10 {
		t.Errorf("Expected short URL length of 10, got %d", len(shortURL))
	}

	for _, char := range shortURL {
		if !isValidCharacter(char) {
			t.Errorf("Generated short URL contains invalid character: %c", char)
		}
	}
}

func isValidCharacter(char rune) bool {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	for _, validChar := range charset {
		if char == validChar {
			return true
		}
	}
	return false
}
