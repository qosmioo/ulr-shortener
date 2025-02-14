package storage

import (
	"testing"
)

type mockRedisStorage struct {
	urls map[string]string
}

func (m *mockRedisStorage) SaveURL(shortURL, originalURL string) error {
	m.urls[shortURL] = originalURL
	return nil
}

func (m *mockRedisStorage) GetURL(shortURL string) (string, error) {
	return m.urls[shortURL], nil
}

func (m *mockRedisStorage) GetShortURLByOriginal(originalURL string) (string, error) {
	for short, original := range m.urls {
		if original == originalURL {
			return short, nil
		}
	}
	return "", ErrURLExists
}

func TestRedisStorage(t *testing.T) {
	mock := &mockRedisStorage{urls: make(map[string]string)}

	tests := []struct {
		shortURL    string
		originalURL string
	}{
		{"short456", "http://example.org"},
		{"short789", "http://example.net"},
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
	}
}
