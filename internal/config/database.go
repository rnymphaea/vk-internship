package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	DBName   string `env:"POSTGRES_DB_NAME,required"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`

	PoolMaxConns          int32         `env:"POSTGRES_POOL_MAX_CONNS" envDefault:"4"`
	PoolMinConns          int32         `env:"POSTGRES_POOL_MAX_CONNS" envDefault:"0"`
	PoolMaxConnLifetime   time.Duration `env:"POSTGRES_POOL_MAX_CONN_LIFETIME" envDefault:"1h"`
	PoolMaxConnIdleTime   time.Duration `env:"POSTGRES_POOL_MAX_CONN_IDLE_TIME" envDefault:"30m"`
	PoolHealthCheckPeriod time.Duration `env:"POSTGRES_POOL_HEALTHCHECK_PERIOD" envDefault:"1m"`
}

func LoadPostgresConfig() (*PostgresConfig, error) {
	var cfg PostgresConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
