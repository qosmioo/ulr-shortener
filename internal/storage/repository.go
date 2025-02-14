package storage

import "fmt"

type Storage interface {
	SaveURL(shortURL, originalURL string) error
	GetURL(shortURL string) (string, error)
	GetShortURLByOriginal(originalURL string) (string, error)
}

var ErrURLExists = fmt.Errorf("URL already exists")
