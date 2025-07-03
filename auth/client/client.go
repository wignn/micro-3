package client

import (
	"context"
	"github.com/wignn/micro-3/auth/genproto"
	"google.golang.org/grpc"
)

type AuthClient struct {
	conn    *grpc.ClientConn
	service genproto.AuthServiceClient
}

func NewClient(url string) (*AuthClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := genproto.NewAuthServiceClient(conn)
	return &AuthClient{conn, c}, nil
}

func (c *AuthClient) Close() {
	c.conn.Close()
}
func (cl *AuthClient) Login(c context.Context, email, password string) (*genproto.PostAuthResponse, error) {
	r, err := cl.service.Login(
		c,
		&genproto.PostAuthRequest{
			Email:    email,
			Password: password,
		},
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (cl *AuthClient) RefreshToken(c context.Context, refreshToken string) (*genproto.BackendToken, error) {
	r, err := cl.service.RefreshToken(
		c,
		&genproto.PostRefreshTokenRequest{
			RefreshToken: refreshToken,
		},
	)

	if err != nil {
		return nil, err
	}
	return r, nil
}
