package service

import (
	"context"

	"github.com/wignn/micro-3/auth/model"
)

type AuthService interface {
	Login(c context.Context,email, password string) (*model.AuthResponse, error)
	Logout(c context.Context, email string) error
	RefreshToken(c context.Context, RefreshToken string) (*model.Token, error)
}

func NewAuthService() AuthService {
	return &
}


func (s *authService) Login(c context.Context, email, password string) (*model.AuthResponse, error) {
	return nil, nil
}
