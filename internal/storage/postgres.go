package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/qosmioo/ulr-shortener/config"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const (
	createTableQuery = `
	CREATE TABLE IF NOT EXISTS urls (
		short_url VARCHAR PRIMARY KEY,
		original_url VARCHAR NOT NULL UNIQUE
	);
	CREATE EXTENSION IF NOT EXISTS pg_trgm;
	CREATE INDEX IF NOT EXISTS idx_original_url_trgm ON urls USING gin (original_url gin_trgm_ops);`

	insertURLQuery                = "INSERT INTO urls (short_url, original_url) VALUES ($1, $2)"
	selectURLQuery                = "SELECT original_url FROM urls WHERE short_url = $1"
	selectShortURLByOriginalQuery = "SELECT short_url FROM urls WHERE original_url = $1"
)

type PostgresStorage struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresStorage(cfg config.PostgresConfig, logger *zap.Logger) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Error("Failed to parse config", zap.Error(err))
		return nil, errors.New("failed to parse config")
	}

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, errors.New("failed to connect to database")
	}

	_, err = db.Exec(context.Background(), createTableQuery)
	if err != nil {
		logger.Error("Failed to create table or index", zap.Error(err))
		return nil, errors.New("failed to create table or index")
	}

	return &PostgresStorage{db: db, logger: logger}, nil
}

func (p *PostgresStorage) SaveURL(shortURL, originalURL string) error {
	_, err := p.db.Exec(context.Background(), insertURLQuery, shortURL, originalURL)
	if err != nil {
		p.logger.Error("Failed to save URL", zap.Error(err))
		return errors.New("failed to save URL")
	}
	p.logger.Info("URL saved", zap.String("shortURL", shortURL), zap.String("originalURL", originalURL))
	return nil
}

func (p *PostgresStorage) GetURL(shortURL string) (string, error) {
	var originalURL string
	err := p.db.QueryRow(context.Background(), selectURLQuery, shortURL).Scan(&originalURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.logger.Warn("URL not found", zap.String("shortURL", shortURL))
			return "", errors.New("URL not found")
		}
		p.logger.Error("Failed to get URL", zap.Error(err))
		return "", errors.New("failed to get URL")
	}
	p.logger.Info("URL retrieved", zap.String("shortURL", shortURL), zap.String("originalURL", originalURL))
	return originalURL, nil
}

func (p *PostgresStorage) GetShortURLByOriginal(originalURL string) (string, error) {
	var shortURL string
	err := p.db.QueryRow(context.Background(), selectShortURLByOriginalQuery, originalURL).Scan(&shortURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.logger.Warn("Original URL not found", zap.String("originalURL", originalURL))
			return "", errors.New("original URL not found")
		}
		p.logger.Error("Failed to get short URL by original", zap.Error(err))
		return "", errors.New("failed to get short URL by original")
	}
	p.logger.Info("Short URL retrieved by original", zap.String("shortURL", shortURL), zap.String("originalURL", originalURL))
	return shortURL, nil
}
