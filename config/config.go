package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		HttpPort string `yaml:"http_port"`
		GrpcPort string `yaml:"grpc_port"`
	} `yaml:"server"`
	Database struct {
		Type     string `yaml:"type"`
		Postgres struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			DBName   string `yaml:"dbname"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
		} `yaml:"postgres"`
		Redis struct {
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
		} `yaml:"redis"`
	} `yaml:"database"`
}

type RedisConfig struct {
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type PostgresConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	if user := os.Getenv("PG_USER"); user != "" {
		cfg.Database.Postgres.User = user
	}
	if password := os.Getenv("PG_PASSWORD"); password != "" {
		cfg.Database.Postgres.Password = password
	}
	if dbName := os.Getenv("PG_DB"); dbName != "" {
		cfg.Database.Postgres.DBName = dbName
	}

	return &cfg, nil
}
