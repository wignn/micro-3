package server

import (
	"context"
	"fmt"
	"github.com/wignn/micro-3/account/genproto"
	"github.com/wignn/micro-3/account/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type grpcServer struct {
	service service.AccountService
	genproto.UnimplementedAccountServiceServer
}

func ListenGRPC(s service.AccountService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	genproto.RegisterAccountServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostAccount(ctx context.Context, req *genproto.PostAccountRequest) (*genproto.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &genproto.PostAccountResponse{
		Account: &genproto.Account{
			Id:   a.ID,
			Name: a.Name,
			Email: a.Email,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, req *genproto.GetAccountRequest) (*genproto.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &genproto.GetAccountResponse{
		Account: &genproto.Account{
			Id:   a.ID,
			Name: a.Name,
			Email: a.Email,
		},
	}, nil
}


func (s *grpcServer) GetAccounts(ctx context.Context, req *genproto.GetAccountsRequest) (*genproto.GetAccountsResponse, error) {
	res, err := s.service.ListAccount(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	var accounts []*genproto.Account
	for _, a := range res {
		accounts = append(accounts, &genproto.Account{
			Id:    a.ID,
			Name:  a.Name,
			Email: a.Email,
		})
	}

	return &genproto.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}
