package service

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/review/model"
	"github.com/wignn/micro-3/review/repository"
)

type ReviewService interface {
	PutReview(c context.Context, productId, acccountId string, rating int, comment string) (*model.Review, error)
	GetReviewById(c context.Context, id string) (*model.Review, error)
	GetReviewByProductAndUser(c context.Context, id string, skip, take uint64) ([]*model.Review, error)
}

type reviewService struct {
	repository repository.ReviewRepository
}

func NewReviewService(r repository.ReviewRepository) ReviewService {
	return &reviewService{r}
}

func (s *reviewService) PutReview(c context.Context, productId, acccountId string, rating int, content string) (*model.Review, error) {
	rev := &model.Review{
		ID:        ksuid.New().String(),
		ProductID: productId,
		AccountID: acccountId,
		Rating:    rating,
		Content:   content,
		CreatedAt: time.Now().UTC(),
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


func (s *reviewService) GetReviewByProductAndUser(c context.Context, id string, skip, take uint64) ([]*model.Review, error) {
	review, err := s.repository.GetReviewByProductAndUser(c, id, skip, take)
	if err != nil {
		return nil, err
	}
	return review, nil
}