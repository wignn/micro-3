package client

import (
	"context"
	"github.com/wignn/micro-3/account/genproto"
	"github.com/wignn/micro-3/account/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountClient struct {
	conn    *grpc.ClientConn
	service genproto.AccountServiceClient
}
func NewClient(url string) (*AccountClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := genproto.NewAccountServiceClient(conn)

	return &AccountClient{conn, c}, nil
}


func (cl *AccountClient) Close() {
	cl.conn.Close()
}

func (cl *AccountClient) PostAccount(c context.Context, name, email, password string) (*model.AccountResponse, error) {
	r, err := cl.service.PostAccount(
		c,
		&genproto.PostAccountRequest{Name: name, Email: email, Password: password},
	)
	if err != nil {
		return nil, err
	}

	return &model.AccountResponse{
		ID:   r.Account.Id,
		Name: r.Account.Name,
		Email: r.Account.Email,
	}, nil
}


func (cl *AccountClient) GetAccount(c context.Context, id string) (*model.AccountResponse, error) {
	r, err := cl.service.GetAccount(
		c,
		&genproto.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}

	return &model.AccountResponse{
		ID:   r.Account.Id,
		Name: r.Account.Name,
		Email: r.Account.Email,
	}, nil
}

func (cl *AccountClient) GetAccounts(c context.Context, skip, take uint64) ([]*model.AccountResponse, error) {
	r, err := cl.service.GetAccounts(
		c,
		&genproto.GetAccountsRequest{Skip: skip, Take: take},
	)
	if err != nil {
		return nil, err
	}

	var accounts []*model.AccountResponse
	for _, a := range r.Accounts {
		accounts = append(accounts, &model.AccountResponse{
			ID:   a.Id,
			Name: a.Name,
			Email: a.Email,
		})
	}

	return accounts, nil
}

func (cl *AccountClient) DeleteAccount(c context.Context, id string) (*genproto.DeleteAccountResponse, error) {
	r, err := cl.service.DeleteAccount(
		c,
		&genproto.DeleteAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}


func (cl *AccountClient) EditAccount(c context.Context, id, name, email, password string) (*model.AccountResponse, error) {
	r, err := cl.service.EditAccount(
		c,
		&genproto.EditAccountRequest{Id: id, Name: name, Email: email, Password: password},
	)
	if err != nil {
		return nil, err
	}

	return &model.AccountResponse{
		ID:   r.Account.Id,
		Name: r.Account.Name,
		Email: r.Account.Email,
	}, nil
}