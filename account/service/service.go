package service

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/account/model"
	"github.com/wignn/micro-3/account/repository"
)

type AccountService  interface {
	PostAccount(c context.Context, name string) (*model.Account, error)
	GetAccount(c context.Context, id string) (*model.Account, error)
	ListAccount(c context.Context, skip uint64, take uint64) ([]*model.Account, error)
}


type accountService struct {
	repository repository.AccountRepository
}

func NewAccountService(r repository.AccountRepository) AccountService {
	return &accountService{r}
}

func (s *accountService) PostAccount(c context.Context, name string) (*model.Account, error) {
	
	a := &model.Account{
		Name: name,
		ID: ksuid.New().String(),
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