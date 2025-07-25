package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"vk-internship/internal/database"
	"vk-internship/internal/database/model"
)

func (p *PostgresDB) CreateUser(user *model.User) (*model.User, error) {
	p.log.Debugf("trying to create user", map[string]interface{}{"info": *user})
	const query = `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, password_hash, created_at
	`

	ctx, cancel := context.WithTimeout(context.TODO(), p.timeout)
	defer cancel()

	var createdUser model.User

	err := p.db.QueryRow(ctx, query, user.Username, user.Password).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Password,
		&createdUser.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return nil, database.ErrUserExists
		}
		return nil, fmt.Errorf("insert user failed: %w", err)
	}

	return &createdUser, nil
}

func (p *PostgresDB) GetUserByUsername(username string) (*model.User, error) {
	const query = `
		SELECT id, username, password_hash, created_at 
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(context.TODO(), p.timeout)
	defer cancel()

	var user model.User
	err := p.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	return &user, nil
}
