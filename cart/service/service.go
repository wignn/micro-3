package service

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/wignn/micro-3/cart/model"
	"github.com/wignn/micro-3/cart/repository"
)

type CartService interface {
	PutCart(ctx context.Context, AccountID, ProductID string, Quantity uint32) error
	GetCartByAccount(ctx context.Context, accountID string) (*model.Cart, error)
	DeleteCart(ctx context.Context, id string) error
	GetCartByID(ctx context.Context, id string) (*model.Cart, error)
	UpdateCartItem(ctx context.Context, cartID, productID string, quantity int) error
}

type cartService struct {
	repository repository.CartRepository
}

func NewCartService(r repository.CartRepository) CartService {
	return &cartService{repository: r}
}

func (s *cartService) PutCart(c context.Context, AccountID, ProductID string, Quantity uint32) error {
	d:= &model.CartPutRequest{
		ProductID: ProductID,
		Quantity:  Quantity,
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		AccountID: AccountID,
	}
	return s.repository.PutCart(c, d)
}



func (s *cartService) GetCartByAccount(c context.Context, accountID string) (*model.Cart, error) {
	return s.repository.GetCartByAccount(c, accountID)
}

func (s *cartService) DeleteCart(c context.Context, id string) error {
	return s.repository.DeleteCart(c, id)
}

func (s *cartService) GetCartByID(c context.Context, id string) (*model.Cart, error) {
	return s.repository.GetCartByID(c, id)
}

func (s *cartService) UpdateCartItem(c context.Context, cartID, productID string, quantity int) error {
	return s.repository.UpdateCartItem(c, cartID, productID, quantity)
}