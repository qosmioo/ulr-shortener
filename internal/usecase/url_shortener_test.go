package usecase

import (
	"testing"

	"github.com/qosmioo/ulr-shortener/internal/storage"
)

type mockStorage struct {
	urls map[string]string
}

func (m *mockStorage) SaveURL(shortURL, originalURL string) error {
	m.urls[shortURL] = originalURL
	return nil
}

func (m *mockStorage) GetURL(shortURL string) (string, error) {
	return m.urls[shortURL], nil
}

func (m *mockStorage) GetShortURLByOriginal(originalURL string) (string, error) {
	for short, original := range m.urls {
		if original == originalURL {
			return short, nil
		}
	}
	return "", storage.ErrURLExists
}

func TestCreateShortURL(t *testing.T) {
	mock := &mockStorage{urls: make(map[string]string)}
	usecase := NewURLShortenerService(mock)

	tests := []struct {
		originalURL string
	}{
		{"http://example.com"},
		{"http://example.org"},
	}

	for _, test := range tests {
		shortURL, err := usecase.CreateShortURL(test.originalURL)
		if err != nil {
			t.Fatalf("Failed to create short URL: %v", err)
		}

		if shortURL == "" {
			t.Error("Expected a short URL, got empty string")
		}

		retrievedURL, err := mock.GetURL(shortURL)
		if err != nil {
			t.Fatalf("Failed to get original URL: %v", err)
		}
		if retrievedURL != test.originalURL {
			t.Errorf("Expected %s, got %s", test.originalURL, retrievedURL)
		}
	}
}
