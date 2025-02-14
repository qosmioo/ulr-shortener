package validator

import (
	"errors"
	"net/url"
	"regexp"
)

var (
	ErrInvalidURLFormat = errors.New("invalid URL format")
	ErrInvalidShortURL  = errors.New("short URL must be exactly 10 characters long and contain only valid characters")
	ErrInvalidURLScheme = errors.New("invalid URL scheme: must be http or https")
	ErrInvalidURLHost   = errors.New("invalid URL host: must contain valid domain format")
)

var (
	validShortURLPattern = regexp.MustCompile(`^[a-zA-Z0-9_]{10}$`)
	validHostPattern     = regexp.MustCompile(`^[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
)

func ValidateShortURL(shortURL string) error {
	if !validShortURLPattern.MatchString(shortURL) {
		return ErrInvalidShortURL
	}
	return nil
}

func ValidateURL(originalURL string) error {
	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return ErrInvalidURLFormat
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidURLScheme
	}
	host := parsedURL.Hostname()
	if !validHostPattern.MatchString(host) {
		return ErrInvalidURLHost
	}
	return nil
}
