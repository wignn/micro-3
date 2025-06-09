package repository

import (
	"context"
	"database/sql"

	"github.com/wignn/micro-3/review/model"
)

type ReviewRepository interface {
	Close()
	PutReview(c context.Context, r *model.Review) error
	GetReviewById(c context.Context, id string) (*model.Review, error)
	ListReviewsByProduct(c context.Context, productId string, skip uint64, take uint64) ([]*model.Review, error)
	ListReviewsByUser(c context.Context, userId string, skip uint64, take uint64) ([]*model.Review, error)
	GetReviewByProductAndUser(c context.Context, productId, userId string) (*model.Review, error)
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
	_, err := r.db.ExecContext(c, "INSERT INTO reviews (id, product_id, user_id, rating, comment, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		rev.ID, rev.ProductID, rev.UserID, rev.Rating, rev.Comment, rev.CreatedAt, rev.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}


func (r *PostgresRepository) GetReviewById(c context.Context, id string) (*model.Review, error) {
	row := r.db.QueryRowContext(c, "SELECT id, product_id, user_id, rating, comment, created_at, updated_at FROM reviews WHERE id = $1", id)
	rev := &model.Review{}
	if err := row.Scan(&rev.ID, &rev.ProductID, &rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rev, nil
}

func (r *PostgresRepository) ListReviewsByProduct(c context.Context, productId string, skip uint64, take uint64) ([]*model.Review, error) {
	rows, err := r.db.QueryContext(c, "SELECT id, product_id, user_id, rating, comment, created_at, updated_at FROM reviews WHERE product_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", productId, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		rev := &model.Review{}
		if err := rows.Scan(&rev.ID, &rev.ProductID, &rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.UpdatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *PostgresRepository) ListReviewsByUser (c context.Context, userId string, skip uint64, take uint64) ([]*model.Review, error) {
	rows, err := r.db.QueryContext(c, "SELECT id, product_id, user_id, rating, comment, created_at, updated_at FROM reviews WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", userId, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		rev := &model.Review{}
		if err := rows.Scan(&rev.ID, &rev.ProductID, &rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.UpdatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}


func (r *PostgresRepository) GetReviewByProductAndUser(c context.Context, productId, userId string) (*model.Review, error) {
	row := r.db.QueryRowContext(c, "SELECT id, product_id, user_id, rating, comment, created_at, updated_at FROM reviews WHERE product_id = $1 AND user_id = $2", productId, userId)
	rev := &model.Review{}
	if err := row.Scan(&rev.ID, &rev.ProductID, &rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rev, nil
}