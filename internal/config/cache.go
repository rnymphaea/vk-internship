package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type RedisConfig struct {
	Addr        string        `env:"REDIS_ADDR"`
	Password    string        `env:"REDIS_PASSWORD,required"`
	User        string        `env:"REDIS_USER,required"`
	DB          int           `env:"REDIS_DB,required"`
	MaxRetries  int           `env:"REDIS_MAX_RETRIES" envDefault:"3"`
	DialTimeout time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"10s"`
	Timeout     time.Duration `env:"REDIS_TIMEOUT" envDefault:"5s"`
}

func LoadRedisConfig() (*RedisConfig, error) {
	var cfg RedisConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
