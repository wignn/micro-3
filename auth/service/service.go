package service

import (
	"context"
	"github.com/wignn/micro-3/auth/model"
	"github.com/wignn/micro-3/auth/repository"
	"github.com/wignn/micro-3/auth/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(c context.Context, l *model.AuthRequest) (*model.AuthResponse, error)
	RefreshToken(c context.Context, RefreshToken string) (*model.Token, error)
}

type authService struct {
	repository repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{repository: r}
}

func (s authService) Login(c context.Context, l *model.AuthRequest) (*model.AuthResponse, error) {
	account, err := s.repository.GetAccount(c, l.Email)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, nil
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(l.Password)); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(account.Email)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		ID: 		account.ID,
		Email:      account.Email,
		BackendToken: *token,
	}, nil
}

func (s authService) RefreshToken(c context.Context, refreshToken string) (*model.Token, error) {
	email, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	newToken, err := utils.GenerateToken(email)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

