package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/qosmioo/ulr-shortener/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisStorage struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisStorage(cfg config.RedisConfig, logger *zap.Logger) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return nil, errors.New("failed to connect to Redis")
	}

	logger.Info("Connected to Redis", zap.String("host", cfg.Host), zap.String("port", cfg.Port))
	return &RedisStorage{client: client, logger: logger}, nil
}

func (r *RedisStorage) SaveURL(shortURL, originalURL string) error {
	existingURL, err := r.client.Get(context.Background(), shortURL).Result()
	if err == nil && existingURL == originalURL {
		r.logger.Warn("URL already exists in Redis", zap.String("shortURL", shortURL))
		return ErrURLExists
	}

	if err := r.client.Set(context.Background(), shortURL, originalURL, 0).Err(); err != nil {
		r.logger.Error("Failed to save URL in Redis", zap.Error(err))
		return errors.New("failed to save URL")
	}
	r.logger.Info("URL saved in Redis", zap.String("shortURL", shortURL), zap.String("originalURL", originalURL))
	return nil
}

func (r *RedisStorage) GetURL(shortURL string) (string, error) {
	originalURL, err := r.client.Get(context.Background(), shortURL).Result()
	if err == redis.Nil {
		r.logger.Warn("URL not found in Redis", zap.String("shortURL", shortURL))
		return "", errors.New("URL not found")
	} else if err != nil {
		r.logger.Error("Failed to get URL from Redis", zap.Error(err))
		return "", errors.New("failed to get URL")
	}
	r.logger.Info("URL retrieved from Redis", zap.String("shortURL", shortURL), zap.String("originalURL", originalURL))
	return originalURL, nil
}

func (r *RedisStorage) GetShortURLByOriginal(originalURL string) (string, error) {
	var shortURL string
	iter := r.client.Scan(context.Background(), 0, "*", 0).Iterator()
	for iter.Next(context.Background()) {
		key := iter.Val()
		val, err := r.client.Get(context.Background(), key).Result()
		if err != nil {
			continue
		}
		if val == originalURL {
			shortURL = key
			break
		}
	}
	if err := iter.Err(); err != nil {
		r.logger.Error("Failed to scan Redis", zap.Error(err))
		return "", errors.New("failed to scan Redis")
	}
	if shortURL == "" {
		r.logger.Warn("Original URL not found in Redis", zap.String("originalURL", originalURL))
		return "", errors.New("original URL not found")
	}
	r.logger.Info("Short URL retrieved by original from Redis", zap.String("shortURL", shortURL), zap.String("originalURL", originalURL))
	return shortURL, nil
}
