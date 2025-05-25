package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/wignn/micro-3/catalog/genproto"
	"github.com/wignn/micro-3/catalog/model"
	"github.com/wignn/micro-3/catalog/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service service.CatalogService
	genproto.UnimplementedCatalogServiceServer
}

func ListenGRPC(s service.CatalogService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	genproto.RegisterCatalogServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *genproto.PostProductRequest) (*genproto.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &genproto.PostProductResponse{Product: &genproto.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *genproto.GetProductRequest) (*genproto.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &genproto.GetProductResponse{
		Product: &genproto.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *genproto.GetProductsRequest) (*genproto.GetProductsResponse, error) {
	var res []*model.Product
	var err error
	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductsByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []*genproto.Product{}
	for _, p := range res {
		products = append(
			products,
			&genproto.Product{
				Id:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		)
	}
	return &genproto.GetProductsResponse{Products: products}, nil
}