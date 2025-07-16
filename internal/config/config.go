package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	Port      string        `env:"PORT,required"`
	JWTSecret string        `env:"JWT_SECRET",required`
	JWTTTL    time.Duration `env:"JWT_TTL" envDefault:"24h"`
}

type StorageConfig struct {
	DBType    string `env:"DB_TYPE,required"`
	CacheType string `env:"CACHE_TYPE,required"`
}

func LoadServerConfig() (*ServerConfig, error) {
	var cfg ServerConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadStorageConfig() (*StorageConfig, error) {
	var cfg StorageConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
