package usecase

import (
	"github.com/qosmioo/ulr-shortener/internal/storage"
	"github.com/qosmioo/ulr-shortener/pkg/shortener"
	"github.com/qosmioo/ulr-shortener/pkg/validator"
)

type URLShortenerService struct {
	store storage.Storage
}

func NewURLShortenerService(store storage.Storage) *URLShortenerService {
	return &URLShortenerService{store: store}
}

func (u *URLShortenerService) CreateShortURL(originalURL string) (string, error) {
	if err := validator.ValidateURL(originalURL); err != nil {
		return "", err
	}

	existingShortURL, err := u.store.GetShortURLByOriginal(originalURL)
	if err == nil {
		return existingShortURL, storage.ErrURLExists
	}

	shortURL := shortener.GenerateShortURL()

	if err := validator.ValidateShortURL(shortURL); err != nil {
		return "", err
	}

	if err := u.store.SaveURL(shortURL, originalURL); err != nil {
		return "", err
	}
	return shortURL, nil
}

func (u *URLShortenerService) GetOriginalURL(shortURL string) (string, error) {
	return u.store.GetURL(shortURL)
}
