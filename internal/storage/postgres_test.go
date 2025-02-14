package storage

import (
	"testing"
)

type mockPostgresStorage struct {
	urls map[string]string
}

func (m *mockPostgresStorage) SaveURL(shortURL, originalURL string) error {
	m.urls[shortURL] = originalURL
	return nil
}

func (m *mockPostgresStorage) GetURL(shortURL string) (string, error) {
	return m.urls[shortURL], nil
}

func (m *mockPostgresStorage) GetShortURLByOriginal(originalURL string) (string, error) {
	for short, original := range m.urls {
		if original == originalURL {
			return short, nil
		}
	}
	return "", ErrURLExists
}

func TestPostgresStorage(t *testing.T) {
	mock := &mockPostgresStorage{urls: make(map[string]string)}

	tests := []struct {
		shortURL    string
		originalURL string
	}{
		{"short123", "http://example.com"},
		{"short456", "http://example.org"},
	}

	for _, test := range tests {
		if err := mock.SaveURL(test.shortURL, test.originalURL); err != nil {
			t.Errorf("Failed to save URL: %v", err)
		}

		retrievedURL, err := mock.GetURL(test.shortURL)
		if err != nil {
			t.Errorf("Failed to get URL: %v", err)
		}
		if retrievedURL != test.originalURL {
			t.Errorf("Expected %s, got %s", test.originalURL, retrievedURL)
		}

		retrievedShortURL, err := mock.GetShortURLByOriginal(test.originalURL)
		if err != nil {
			t.Errorf("Failed to get short URL by original: %v", err)
		}
		if retrievedShortURL != test.shortURL {
			t.Errorf("Expected %s, got %s", test.shortURL, retrievedShortURL)
		}
	}
}
