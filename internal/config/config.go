package config

import (
	"errors"
	"os"
)

type Config struct {
	DbURL        string
	LoggingLevel string
}

func Get() (*Config, error) {
	// Get environment variables
	database := os.Getenv("DATABASE_URL")
	if database == "" {
		return nil, errors.New("missing DATABASE_URL env variable")
	}

	loggingLevel := os.Getenv("LOG_LEVEL")
	if loggingLevel == "" {
		return nil, errors.New("missing LOG_LEVEL env variable")
	}

	return &Config{
		DbURL:        database,
		LoggingLevel: loggingLevel,
	}, nil
}
