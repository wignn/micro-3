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

func (s *grpcServer) PostProduct(c context.Context, r *genproto.PostProductRequest) (*genproto.PostProductResponse, error) {
	p, err := s.service.PostProduct(c, r.Name, r.Description, r.Price, r.Image)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &genproto.PostProductResponse{Product: &genproto.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Image:       p.Image,
	}}, nil
}

func (s *grpcServer) GetProduct(c context.Context, r *genproto.GetProductRequest) (*genproto.GetProductResponse, error) {
	p, err := s.service.GetProduct(c, r.Id)
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
			Image:       p.Image,
		},
	}, nil
}

func (s *grpcServer) GetProducts(c context.Context, r *genproto.GetProductsRequest) (*genproto.GetProductsResponse, error) {
	var res []*model.Product
	var err error

	if r.Query != "" {
		res, err = s.service.SearchProducts(c, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductsByIDs(c, r.Ids)
	} else {
		res, err = s.service.GetProducts(c, r.Skip, r.Take)
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
				Image:       p.Image,
			},
		)
	}
	return &genproto.GetProductsResponse{Products: products}, nil
}

func (s *grpcServer) EditProduct(c context.Context, r *genproto.EditProductRequest) (*genproto.PostProductResponse, error) {
	p, err := s.service.EditProduct(c, r.Id, r.Name, r.Description, r.Price, r.Image)
	if err != nil {
		log.Println("failed to edit product:", err)
		return nil, err 
}
	return &genproto.PostProductResponse{Product: &genproto.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Image:       p.Image,
	}}, nil
}

func (s *grpcServer) DeleteProduct(c context.Context, r *genproto.DeleteProductRequest) (*genproto.DeleteProductResponse, error) {
	err := s.service.DeleteProduct(c, r.Id)
	if err != nil {
		log.Printf("failed to delete product with ID %s: %v\n", r.Id, err)
		return nil,  err
	}

	return &genproto.DeleteProductResponse{
		DeletedID: r.Id,
		Message:   fmt.Sprintf("Product with ID %s deleted successfully", r.Id),
		Success:   true,
	}, nil
}
