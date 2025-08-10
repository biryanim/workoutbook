package config

import (
	"github.com/joho/godotenv"
	"time"
)

type PGConfig interface {
	DSN() string
}

type HTTPConfig interface {
	Address() string
}

type JWTConfig interface {
	TokenSecret() []byte
	TokenExpiration() time.Duration
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
