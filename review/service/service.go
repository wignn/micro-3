package service

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/review/model"
	"github.com/wignn/micro-3/review/repository"
)

type ReviewService interface {
	PutReview(c context.Context, productId, userId string, rating int, comment string) (*model.Review, error)
	GetReviewById(c context.Context, id string) (*model.Review, error)
	ListReviewsByProduct(c context.Context, productId string, skip uint64, take uint64) ([]*model.Review, error)
	ListReviewsByUser(c context.Context, userId string, skip uint64, take uint64) ([]*model.Review, error)
	GetReviewByProductAndUser(c context.Context, productId, userId string) (*model.Review, error)
}

type reviewService struct {
	repository repository.ReviewRepository
}

func NewReviewService(r repository.ReviewRepository) ReviewService {
	return &reviewService{r}
}

func (s *reviewService) PutReview(c context.Context, productId, userId string, rating int, comment string) (*model.Review, error) {
	rev := &model.Review{
		ID:        ksuid.New().String(),
		ProductID: productId,
		UserID:    userId,
		Rating:    rating,
		Comment:   comment,
	}

	if err := s.repository.PutReview(c, rev); err != nil {
		return nil, err
	}

	return rev, nil
}

func (s *reviewService) GetReviewById(c context.Context, id string) (*model.Review, error) {
	rev, err := s.repository.GetReviewById(c, id)
	if err != nil {
		return nil, err
	}
	return rev, nil
}

func (s *reviewService) ListReviewsByProduct(c context.Context, productId string, skip uint64, take uint64) ([]*model.Review, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}

	reviews, err := s.repository.ListReviewsByProduct(c, productId, skip, take)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *reviewService) ListReviewsByUser(c context.Context, userId string, skip uint64, take uint64) ([]*model.Review, error) {
	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}

	reviews, err := s.repository.ListReviewsByUser(c, userId, skip, take)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *reviewService) GetReviewByProductAndUser(c context.Context, productId, userId string) (*model.Review, error) {
	review, err := s.repository.GetReviewByProductAndUser(c, productId, userId)
	if err != nil {
		return nil, err
	}
	return review, nil
}