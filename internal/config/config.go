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

type LoggerConfig struct {
	Type   string `env:"LOGGER_TYPE" envDefault:"zerolog"`
	Level  string `env:"LOGGER_LEVEL" envDefault:"info"`
	Pretty bool   `env:"LOGGER_PRETTY" envDefault:"false"`
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

func LoadLoggerConfig() (*LoggerConfig, error) {
	var cfg LoggerConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
