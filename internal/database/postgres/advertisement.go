package postgres

import (
	"context"
	"fmt"

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

	return &model.Advertisement{
		ID:          createdAd.ID,
		AuthorID:    createdAd.AuthorID,
		Caption:     createdAd.Caption,
		Description: createdAd.Description,
		ImageURL:    createdAd.ImageURL,
		Price:       createdAd.Price,
		CreatedAt:   createdAd.CreatedAt,
	}, nil
}
