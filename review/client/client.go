package client

import (
	"context"
	"time"

	"github.com/wignn/micro-3/review/genproto"
	"github.com/wignn/micro-3/review/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReviewClient struct {
	conn    *grpc.ClientConn
	service genproto.ReviewServiceClient
}

func NewClient(url string) (*ReviewClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := genproto.NewReviewServiceClient(conn)

	return &ReviewClient{conn, c}, nil
}

func (cl *ReviewClient) Close() {
	cl.conn.Close()
}

func (cl *ReviewClient) PostReview(c context.Context, productId, accountId, content string, rating int32) (*model.Review, error) {
	r, err := cl.service.PostReview(
		c, &genproto.PostReviewRequest{
			ProductId: productId,
			AccountId: accountId,
			Content:   content,
			Rating:    rating},
	)
	if err != nil {
		return nil, err
	}
	newDate := time.Time{}
	newDate.UnmarshalBinary(r.Review.CreatedAt)

	return &model.Review{
		ID:        r.Review.Id,
		ProductID: r.Review.ProductId,
		Rating:    int(r.Review.Rating),
		Content:   r.Review.Content,
		AccountID: r.Review.AccountId,
		CreatedAt: newDate,
	}, nil
}

func (cl *ReviewClient) GetReview(c context.Context, id string) (*genproto.Review, error) {
	r, err := cl.service.GetReview(c, &genproto.ReviewIdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &genproto.Review{
		Id:        r.Id,
		ProductId: r.ProductId,
		AccountId: r.AccountId,
		Content:  r.Content,
		Rating:    r.Rating,
		CreatedAt: r.CreatedAt,
	}, nil
}

func (cl *ReviewClient) GetReviews(c context.Context, id string, skip, take uint64) ([]*genproto.Review, error) {
	r, err := cl.service.GetReviewByProductAndUser(c, &genproto.ProductIdRequest{
		Id:   id,
		Skip: skip,
		Take: take})
	if err != nil {
		return nil, err
	}

	var reviews []*genproto.Review

	for _, rev := range r.Reviews {
		reviews = append(reviews, &genproto.Review{
			Id:        rev.Id,
			ProductId: rev.ProductId,
			AccountId: rev.AccountId,
			Rating:    rev.Rating,
			Content:   rev.Content,
			CreatedAt: rev.CreatedAt,
		})
	}
	return reviews, nil
}
