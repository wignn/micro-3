package client

import (
	"context"
	"log"
	"time"

	"github.com/wignn/micro-3/order/genproto"
	"github.com/wignn/micro-3/order/model"
	"google.golang.org/grpc"
)

type OrderClient struct {
	conn    *grpc.ClientConn
	service genproto.OrderServiceClient
}

func NewClient(url string) (*OrderClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := genproto.NewOrderServiceClient(conn)
	return &OrderClient{conn, c}, nil
}

func (c *OrderClient) Close() {
	c.conn.Close()
}

func (c *OrderClient) PostOrder(
	ctx context.Context,
	accountID string,
	products []*model.OrderedProduct,
) (*model.Order, error) {
	protoProducts := []*genproto.PostOrderRequest_OrderProduct{}
	for _, p := range products {
		protoProducts = append(protoProducts, &genproto.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}
	r, err := c.service.PostOrder(
		ctx,
		&genproto.PostOrderRequest{
			AccountId: accountID,
			Products:  protoProducts,
		},
	)
	if err != nil {
		return nil, err
	}

	// Create response order
	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)

	// Convert []*model.OrderedProduct to []model.OrderedProduct
	orderProducts := make([]model.OrderedProduct, len(products))
	for i, p := range products {
		orderProducts[i] = *p
	}
	
	return &model.Order{
		ID:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountID:  newOrder.AccountId,
		Products:   orderProducts,
	}, nil
}

func (cl *OrderClient) GetOrdersForAccount(c context.Context, accountID string) ([]model.Order, error) {
	r, err := cl.service.GetOrdersForAccount(c, &genproto.GetOrdersForAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Create response orders
	orders := []model.Order{}
	for _, orderProto := range r.Orders {
		newOrder := model.Order{
			ID:         orderProto.Id,
			TotalPrice: orderProto.TotalPrice,
			AccountID:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)

		products := []model.OrderedProduct{}
		for _, p := range orderProto.Products {
			products = append(products, model.OrderedProduct{
				ID:          p.Id,
				Quantity:    p.Quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
	
			})
		}
		newOrder.Products = products

		orders = append(orders, newOrder)
	}
	return orders, nil
}

