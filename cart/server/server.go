package server

import (
	"context"
	"fmt"
	"net"

	account "github.com/wignn/micro-3/account/client"
	"github.com/wignn/micro-3/cart/genproto"
	"github.com/wignn/micro-3/cart/service"
	catalog "github.com/wignn/micro-3/catalog/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       service.CartService
	accountClient *account.AccountClient
	catalogClient *catalog.CatalogClient
	genproto.UnimplementedCartServiceServer
}

func ListenGRPC(s service.CartService, catalogURL, accountURL string, port int) error {
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		catalogClient.Close()
		return err
	}

	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		catalogClient.Close()
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		catalogClient.Close()
		accountClient.Close()
		return err
	}
	serv := grpc.NewServer()
	genproto.RegisterCartServiceServer(serv, &grpcServer{
		service:       s,
		catalogClient: catalogClient,
		accountClient: accountClient,
	})

	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostCart(ctx context.Context, req *genproto.PostCartRequest) (*genproto.PostCartResponse, error) {
	err := s.service.PutCart(ctx, req.AccountId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	p, err := s.catalogClient.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	a, err := s.accountClient.GetAccount(ctx, req.AccountId)

	if a == nil {
		return nil, fmt.Errorf("account with ID %s not found", req.AccountId)

	}

	if err != nil {
		return nil, err
	}
	return &genproto.PostCartResponse{
		CartProduct: &genproto.CartProduct{
			Id:          req.ProductId,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    req.Quantity,
		},
	}, nil
}
