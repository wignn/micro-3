package client

import (
	"context"

	"github.com/wignn/micro-3/cart/genproto"
	"github.com/wignn/micro-3/cart/model"
	"google.golang.org/grpc"
)

type CartClient struct {
	conn    *grpc.ClientConn
	service genproto.CartServiceClient
}

func NewClient(url string) (*CartClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := genproto.NewCartServiceClient(conn)
	return &CartClient{conn: conn, service: c}, nil
}

func (c *CartClient) Close() {
	c.conn.Close()
}

func (cl *CartClient) PostCart(c context.Context, productID string, accountID string, quantity uint32) (*model.CartPutResponse, error) {
	r, err := cl.service.PostCart(c, &genproto.PostCartRequest{
		ProductId: productID,
		AccountId: accountID,
		Quantity:  quantity,
	})
	if err != nil {
		return nil, err
	}
	return &model.CartPutResponse{
		ID:          r.CartProduct.Id,
		Quantity:    r.CartProduct.Quantity,
		Name:        r.CartProduct.Name,
		Price:       r.CartProduct.Price,
		Description: r.CartProduct.Description,
	}, nil
}
