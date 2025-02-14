package validator

import (
	"errors"
	"net/url"
	"regexp"
)

var (
	ErrInvalidURLFormat = errors.New("invalid URL format")
	ErrInvalidShortURL  = errors.New("short URL must be exactly 10 characters long and contain only valid characters")
)

var validShortURLPattern = regexp.MustCompile(`^[a-zA-Z0-9_]{10}$`)

func ValidateShortURL(shortURL string) error {
	if !validShortURLPattern.MatchString(shortURL) {
		return ErrInvalidShortURL
	}
	return nil
}

func ValidateURL(originalURL string) error {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return ErrInvalidURLFormat
	}
	return nil
}
