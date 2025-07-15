package server

import (
	"context"
	"fmt"
	"net"
	"github.com/wignn/micro-3/review/genproto"
	"github.com/wignn/micro-3/review/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service service.ReviewService
	genproto.UnimplementedReviewServiceServer
}

func ListenGRPC(s service.ReviewService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	genproto.RegisterReviewServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostReview(ctx context.Context, req *genproto.PostReviewRequest) (*genproto.ReviewResponse, error) {
	rev, err := s.service.PutReview(ctx, req.ProductId, req.AccountId, int(req.Rating), req.Content)
	if err != nil {
		return nil, err
	}
	return &genproto.ReviewResponse{Review: &genproto.Review{
		Id:        rev.ID,
		ProductId: rev.ProductID,
		AccountId: rev.AccountID,
		Rating:    int32(rev.Rating),
		Content:   rev.Content,
	},
	}, nil
}

func (s *grpcServer) GetReview(ctx context.Context, req *genproto.ReviewIdRequest) (*genproto.Review, error) {
	rev, err := s.service.GetReviewById(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	if rev == nil {
		return nil, fmt.Errorf("review not found")
	}

	createdAtBytes, _ := rev.CreatedAt.MarshalBinary()

	return &genproto.Review{
		Id:        rev.ID,
		ProductId: rev.ProductID,
		AccountId: rev.AccountID,
		Rating:    int32(rev.Rating),
		Content:   rev.Content,
		CreatedAt: createdAtBytes,
	}, nil
}

func (s *grpcServer) GetReviewByProductAndUser(
	ctx context.Context,
	req *genproto.ProductIdRequest,
) (*genproto.ReviewListResponse, error) {

	rev, err := s.service.GetReviewByProductAndUser(ctx, req.Id, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	if len(rev) == 0 {
		return nil, fmt.Errorf("no reviews found")
	}

	var reviews []*genproto.Review
	for _, r := range rev {
		createdAtBytes, _ := r.CreatedAt.MarshalBinary()

		reviews = append(reviews, &genproto.Review{
			Id:        r.ID,
			ProductId: r.ProductID,
			AccountId: r.AccountID,
			Rating:    int32(r.Rating),
			Content:   r.Content,
			CreatedAt: createdAtBytes,
		})
	}
	return &genproto.ReviewListResponse{Reviews: reviews}, nil
}
