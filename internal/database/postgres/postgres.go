package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"vk-internship/internal/config"
	"vk-internship/internal/logger"
)

const (
	uniqueViolationCode           = "23505"
	invalidTextRepresentationCode = "22P02"
)

type PostgresDB struct {
	db      *pgxpool.Pool
	log     logger.Logger
	timeout time.Duration
}

func New(cfg *config.PostgresConfig, log logger.Logger) (*PostgresDB, error) {
	log.Debug("creating new postgres pool")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	pool, err := pgxpool.New(context.TODO(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	p := &PostgresDB{
		db:      pool,
		log:     log.Component("postgres"),
		timeout: cfg.Timeout,
	}

	if err := p.Ping(context.TODO()); err != nil {
		return nil, err
	}

	p.log.Info("connected to postgres")

	return p, nil
}

func (p *PostgresDB) Ping(ctx context.Context) error {
	p.log.Debug("ping postgres")
	return p.db.Ping(ctx)
}

func (p *PostgresDB) Close() {
	p.db.Close()
}
