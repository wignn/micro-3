package server

import (
	"context"
	"fmt"
	"net"

	"github.com/wignn/micro-3/auth/genproto"
	"github.com/wignn/micro-3/auth/model"
	"github.com/wignn/micro-3/auth/service"
	"google.golang.org/grpc"
)

type grpcServer struct {
	service service.AuthService
	genproto.UnimplementedAuthServiceServer
}

func ListenGRPC(s service.AuthService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer lis.Close()

	serv := grpc.NewServer()
	genproto.RegisterAuthServiceServer(serv, &grpcServer{
		service: s,
	})

	return serv.Serve(lis)
}
func (s *grpcServer) Login(c context.Context, r *genproto.PostAuthRequest) (*genproto.PostAuthResponse, error) {
	authRequest := &model.AuthRequest{
		Email:    r.Email,
		Password: r.Password,
	}

	user, err := s.service.Login(c, authRequest)
	if err != nil {
		return nil, err
	}

	return &genproto.PostAuthResponse{
		Auth: &genproto.Auth{
			Id:    user.ID,
			Email: user.Email,
			Token: &genproto.BackendToken{
				AccessToken:  user.BackendToken.AccessToken,
				RefreshToken: user.BackendToken.RefreshToken,
				ExpiresAt:    user.BackendToken.ExpiresAt,
			},
		},
	}, nil
}

func (s *grpcServer) RefreshToken(c context.Context, r *genproto.PostRefreshTokenRequest) (*genproto.BackendToken, error) {
	newToken, err := s.service.RefreshToken(c, r.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &genproto.BackendToken{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresAt:    newToken.ExpiresAt,
	}, nil
}
