package repository

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/wignn/micro-3/review/model"
)

type ReviewRepository interface {
	Close()
	PutReview(c context.Context, r *model.Review) error
	GetReviewById(c context.Context, id string) (*model.Review, error)
	GetReviewByProductAndUser(c context.Context, id string, skip, take uint64) ([]*model.Review, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Close() {
	if err := r.db.Close(); err != nil {
		panic(err)
	}
}

func (r *PostgresRepository) PutReview(c context.Context, rev *model.Review) error {
	_, err := r.db.ExecContext(c, "INSERT INTO reviews (id, product_id, account_id, rating, content, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		rev.ID, rev.ProductID, rev.AccountID, rev.Rating, rev.Content, rev.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetReviewById(c context.Context, id string) (*model.Review, error) {
	row := r.db.QueryRowContext(c, "SELECT id, product_id, account_id, rating, content, created_at FROM reviews WHERE id = $1", id)
	rev := &model.Review{}
	if err := row.Scan(&rev.ID, &rev.ProductID, &rev.AccountID, &rev.Rating, &rev.Content, &rev.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rev, nil
}

func (r *PostgresRepository) GetReviewByProductAndUser(c context.Context, id string, skip, take uint64) ([]*model.Review, error) {

	rows, err := r.db.QueryContext(c, `
		SELECT id, product_id, account_id, rating, content, created_at
		FROM reviews
		WHERE product_id = $1 OR account_id = $2
		LIMIT $3 OFFSET $4
	`, id, id, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		review := &model.Review{}
		if err := rows.Scan(&review.ID, &review.ProductID, &review.AccountID, &review.Rating, &review.Content, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
