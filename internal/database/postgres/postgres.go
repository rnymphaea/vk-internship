package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"vk-internship/internal/config"
)

type PostgresDB struct {
	db *pgxpool.Pool
}

func New(cfg *config.PostgresConfig) (*PostgresDB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &PostgresDB{db: pool}, nil
}

func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.db.Ping(ctx)
}
