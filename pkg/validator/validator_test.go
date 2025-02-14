package validator

import (
	"testing"
)

func TestValidateShortURL(t *testing.T) {
	tests := []struct {
		shortURL string
		valid    bool
	}{
		{"abcdefghij", true},
		{"1234567890", true},
		{"abcde12345", true},
		{"abc", false},
		{"abcdefghijk", false},
		{"abc$%^&*()", false},
	}

	for _, test := range tests {
		err := ValidateShortURL(test.shortURL)
		if (err == nil) != test.valid {
			t.Errorf("Expected validity for %s to be %v, got %v", test.shortURL, test.valid, err == nil)
		}
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url   string
		valid bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", true},
		{"invalid-url", false},
		{"", false},
	}

	for _, test := range tests {
		err := ValidateURL(test.url)
		if (err == nil) != test.valid {
			t.Errorf("Expected validity for %s to be %v, got %v", test.url, test.valid, err == nil)
		}
	}
}
