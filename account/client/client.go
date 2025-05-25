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

func (cl *AccountClient) PostAccount(c context.Context, name string) (*model.Account, error) {
	r, err := cl.service.PostAccount(
		c,
		&genproto.PostAccountRequest{Name: name},
	)
	if err != nil {
		return nil, err
	}

	return &model.Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}


func (cl *AccountClient) GetAccount(c context.Context, id string) (*model.Account, error) {
	r, err := cl.service.GetAccount(
		c,
		&genproto.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}

	return &model.Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (cl *AccountClient) GetAccounts(c context.Context, skip, take uint64) ([]*model.Account, error) {
	r, err := cl.service.GetAccounts(
		c,
		&genproto.GetAccountsRequest{Skip: skip, Take: take},
	)
	if err != nil {
		return nil, err
	}

	var accounts []*model.Account
	for _, a := range r.Accounts {
		accounts = append(accounts, &model.Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}

	return accounts, nil
}