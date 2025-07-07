package service

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/account/model"
	"github.com/wignn/micro-3/account/repository"
	"golang.org/x/crypto/bcrypt"
)

type AccountService  interface {
	PostAccount(c context.Context, name, email, password string) (*model.Account, error)
	GetAccount(c context.Context, id string) (*model.Account, error)
	ListAccount(c context.Context, skip uint64, take uint64) ([]*model.Account, error)
	DeleteAccount(c context.Context, id string) error
	EditAccount(c context.Context, id, name, email, password string) (*model.Account, error)
}


type accountService struct {
	repository repository.AccountRepository
}

func NewAccountService(r repository.AccountRepository) AccountService {
	return &accountService{r}
}

func (s *accountService) PostAccount(c context.Context, name, email, password string) (*model.Account, error) {
	
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	a := &model.Account{
		Name: name,
		ID: ksuid.New().String(),
		Email: email,
		Password: string(hashPassword),
	}

	if err := s.repository.PutAccount(c, a); err != nil {
		return nil, err
	}

	return a, nil
}


func (s *accountService) GetAccount(c context.Context, id string) (*model.Account, error) {
	a, err := s.repository.GetAccountById(c, id)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *accountService) ListAccount(c context.Context, skip uint64, take uint64) ([]*model.Account, error) {

	if take > 100 || (take == 0 && skip == 0) {
		take = 100
	}

	accounts, err := s.repository.ListAccount(c, skip, take)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *accountService) DeleteAccount(c context.Context, id string) error {
	if id == "" {
		return repository.ErrNotFound
	}

	return s.repository.DeleteAccount(c, id)
}

func (s *accountService) EditAccount(c context.Context, id, name, email, password string) (*model.Account, error) {
	if id == "" {
		return nil, repository.ErrNotFound
	}
	
 	r, err := s.repository.EditAccount(c , &model.Account{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}