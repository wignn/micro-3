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
	rev, err := s.service.PutReview(ctx, req.ProductId, req.UserId, int(req.Rating), req.Content)
	if err != nil {
		return nil, err
	}
	return &genproto.ReviewResponse{Review: &genproto.Review{
		Id:        rev.ID,
		ProductId: rev.ProductID,
		UserId:    rev.UserID,
		Rating:    int32(rev.Rating),
		Content:   rev.Comment,
		CreatedAt: []byte(rev.CreatedAt),
		UpdatedAt: []byte(rev.UpdatedAt),
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
	return &genproto.Review{
		Id:        rev.ID,
		ProductId: rev.ProductID,
		UserId:    rev.UserID,
		Rating:    int32(rev.Rating),
		Content:   rev.Comment,
		CreatedAt: []byte(rev.CreatedAt),
		UpdatedAt: []byte(rev.UpdatedAt),
	}, nil
}

func (s *grpcServer) GetReviewsByProduct(ctx context.Context, req *genproto.ProductIdRequest) (*genproto.ReviewListResponse, error) {
	revs, err := s.service.ListReviewsByProduct(ctx, req.ProductId, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	var protoRevs []*genproto.Review
	for _, rev := range revs {
		protoRevs = append(protoRevs, &genproto.Review{
			Id:        rev.ID,
			ProductId: rev.ProductID,
			UserId:    rev.UserID,
			Rating:    int32(rev.Rating),
			Content:   rev.Comment,
			CreatedAt: []byte(rev.CreatedAt),
			UpdatedAt: []byte(rev.UpdatedAt),
		})
	}
	return &genproto.ReviewListResponse{Reviews: protoRevs}, nil
}

func (s *grpcServer) GetReviewsByUser(ctx context.Context, req *genproto.UserIdRequest) (*genproto.ReviewListResponse, error) {
	revs, err := s.service.ListReviewsByUser(ctx, req.UserId, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	var protoRevs []*genproto.Review
	for _, rev := range revs {
		protoRevs = append(protoRevs, &genproto.Review{
			Id:        rev.ID,
			ProductId: rev.ProductID,
			UserId:    rev.UserID,
			Rating:    int32(rev.Rating),
			Content:   rev.Comment,
			CreatedAt: []byte(rev.CreatedAt),
			UpdatedAt: []byte(rev.UpdatedAt),
		})
	}
	return &genproto.ReviewListResponse{Reviews: protoRevs}, nil
}

func (s *grpcServer) GetReviewByProductAndUser(ctx context.Context, req *genproto.ProductUserRequest) (*genproto.ReviewResponse, error) {
	rev, err := s.service.GetReviewByProductAndUser(ctx, req.ProductId, req.UserId)
	if err != nil {
		return nil, err
	}
	if rev == nil {
		return nil, fmt.Errorf("review not found")
	}
	return &genproto.ReviewResponse{Review: &genproto.Review{
		Id:        rev.ID,
		ProductId: rev.ProductID,
		UserId:    rev.UserID,
		Rating:    int32(rev.Rating),
		Content:   rev.Comment,
		CreatedAt: []byte(rev.CreatedAt),
		UpdatedAt: []byte(rev.UpdatedAt),
	}}, nil
}
