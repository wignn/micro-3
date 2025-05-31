package service

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/order/model"
	"github.com/wignn/micro-3/order/repository"
)

type OrderService interface {
	PostOrder(c context.Context, accountID string, products []model.OrderedProduct) (*model.Order, error)
	GetOrdersForAccount(c context.Context, accountID string) ([]*model.Order, error)
}


type orderService struct {
	repository repository.OrderRepository
}


func NewOrderService(r repository.OrderRepository) OrderService {
	return &orderService{repository: r}
}



func (s orderService) PostOrder(
	ctx context.Context,
	accountID string,
	products []model.OrderedProduct,
) (*model.Order, error) {
	o := &model.Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}
	
	// Calculate total price
	o.TotalPrice = 0.0
	for _, p := range products {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}
	err := s.repository.PutOrder(ctx, o)
	
	if err != nil {
		return nil, err
	}
	
	return o, nil
}

func (s orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]*model.Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountID)
}