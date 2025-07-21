package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"vk-internship/internal/database"
	"vk-internship/internal/database/model"
)

func (p *PostgresDB) CreateAd(ad *model.Advertisement) (*model.Advertisement, error) {
	const query = `
		INSERT INTO advertisements (author_id, caption, description, image_url, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, author_id, caption, description, image_url, price, created_at
	`

	ctx, cancel := context.WithTimeout(context.TODO(), p.timeout)
	defer cancel()

	var createdAd model.Advertisement
	err := p.db.QueryRow(ctx, query,
		ad.AuthorID,
		ad.Caption,
		ad.Description,
		ad.ImageURL,
		ad.Price,
	).Scan(&createdAd.ID,
		&createdAd.AuthorID,
		&createdAd.Caption,
		&createdAd.Description,
		&createdAd.ImageURL,
		&createdAd.Price,
		&createdAd.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert ad failed: %w", err)
	}

	return &createdAd, nil
}

func (p *PostgresDB) GetAds(ctx context.Context, sortBy, order string, minPrice, maxPrice *int, page, pageSize int) ([]*model.Advertisement, int, error) {
	var params []interface{}
	conditions := []string{"1=1"}

	if minPrice != nil {
		conditions = append(conditions, fmt.Sprintf("a.price >= $%d", len(params)+1))
		params = append(params, *minPrice)
	}

	if maxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("a.price <= $%d", len(params)+1))
		params = append(params, *maxPrice)
	}

	whereClause := strings.Join(conditions, " AND ")

	query := fmt.Sprintf(`
        SELECT 
            a.id, 
            a.author_id, 
            u.username as author_username, 
            a.caption, 
            a.description, 
            a.image_url, 
            a.price, 
            a.created_at,
            COUNT(*) OVER() AS total_count
        FROM advertisements a
        JOIN users u ON a.author_id = u.id
        WHERE %s`, whereClause)

	validSortFields := map[string]bool{"created_at": true, "price": true}
	if _, ok := validSortFields[sortBy]; !ok {
		sortBy = "created_at"
	}
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY a.%s %s", sortBy, order)

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(params)+1, len(params)+2)
	params = append(params, offset, pageSize)

	rows, err := p.db.Query(ctx, query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var ads []*model.Advertisement
	totalCount := 0

	for rows.Next() {
		var ad model.Advertisement

		err := rows.Scan(
			&ad.ID,
			&ad.AuthorID,
			&ad.AuthorUsername,
			&ad.Caption,
			&ad.Description,
			&ad.ImageURL,
			&ad.Price,
			&ad.CreatedAt,
			&totalCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}

		ads = append(ads, &ad)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return ads, totalCount, nil
}

func (p *PostgresDB) GetAd(ctx context.Context, id string) (*model.Advertisement, error) {
	p.log.Debugf("get advertisement", map[string]interface{}{"ad_id": id})

	const query = `
        SELECT 
            a.id, 
            a.author_id, 
            u.username as author_username, 
            a.caption, 
            a.description, 
            a.image_url, 
            a.price, 
            a.created_at
        FROM advertisements a
        JOIN users u ON a.author_id = u.id
        WHERE a.id = $1
    `

	var ad model.Advertisement
	err := p.db.QueryRow(ctx, query, id).Scan(
		&ad.ID,
		&ad.AuthorID,
		&ad.AuthorUsername,
		&ad.Caption,
		&ad.Description,
		&ad.ImageURL,
		&ad.Price,
		&ad.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if (errors.As(err, &pgErr) && pgErr.Code == invalidTextRepresentationCode) || errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrAdNotFound
		}

		return nil, fmt.Errorf("failed to get ad: %w", err)
	}

	return &ad, nil
}

func (p *PostgresDB) DeleteAd(ctx context.Context, id, authorID string) error {
	p.log.Debugf("delete ad", map[string]interface{}{"ad_id": id, "author_id": authorID})

	const query = `DELETE FROM advertisements WHERE id = $1 AND author_id = $2`

	result, err := p.db.Exec(ctx, query, id, authorID)
	if err != nil {
		return fmt.Errorf("failed to delete ad: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return database.ErrAdNotFoundOrNotOwnedByUser
	}

	return nil
}

func (p *PostgresDB) UpdateAd(ctx context.Context, ad *model.Advertisement) (*model.Advertisement, error) {
	const query = `
        UPDATE advertisements
        SET 
            caption = $1,
            description = $2,
            image_url = $3,
            price = $4,
						updated_at = $5
        WHERE id = $6 AND author_id = $7
        RETURNING id, author_id, caption, description, image_url, price, created_at, updated_at
    `

	var updatedAd model.Advertisement
	err := p.db.QueryRow(ctx, query,
		ad.Caption,
		ad.Description,
		ad.ImageURL,
		ad.Price,
		time.Now(),
		ad.ID,
		ad.AuthorID,
	).Scan(
		&updatedAd.ID,
		&updatedAd.AuthorID,
		&updatedAd.Caption,
		&updatedAd.Description,
		&updatedAd.ImageURL,
		&updatedAd.Price,
		&updatedAd.CreatedAt,
		&updatedAd.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update ad: %w", err)
	}

	return &updatedAd, nil
}
