package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/wignn/micro-3/account/genproto"
	"github.com/wignn/micro-3/account/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func (s *grpcServer) PostAccount(c context.Context, req *genproto.PostAccountRequest) (*genproto.PostAccountResponse, error) {
	a, err := s.service.PostAccount(c, req.Name, req.Email, req.Password)
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

func (s *grpcServer) GetAccount(c context.Context, req *genproto.GetAccountRequest) (*genproto.GetAccountResponse, error) {
	a, err := s.service.GetAccount(c, req.Id)
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


func (s *grpcServer) GetAccounts(c context.Context, req *genproto.GetAccountsRequest) (*genproto.GetAccountsResponse, error) {
	res, err := s.service.ListAccount(c, req.Skip, req.Take)
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

func (s *grpcServer) DeleteAccount(c context.Context, req *genproto.DeleteAccountRequest) (*genproto.DeleteAccountResponse, error) {
	if err := s.service.DeleteAccount(c, req.Id); err != nil {
		return nil, err
	}

	return &genproto.DeleteAccountResponse{
		Message:   "Account deleted successfully",
		Success:   true,
		DeletedID: req.Id,
	}, nil
}


func (s *grpcServer) EditAccount(c context.Context, req *genproto.EditAccountRequest) (*genproto.EditAccountResponse, error) {
	a, err := s.service.EditAccount(c, req.Id, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	log.Printf("Account with ID %s updated successfully", a.ID)
	return &genproto.EditAccountResponse{
		Message: "Account updated successfully",
		Success: true,
		Account: &genproto.Account{
			Id:    a.ID,
			Name:  a.Name,
			Email: a.Email,
		},
	}, nil
}