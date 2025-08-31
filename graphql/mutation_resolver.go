package main

import (
	"context"
	"errors"
	"fmt"
	productModel "github.com/wignn/micro-3/order/model"
	"time"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *GraphQLServer
}

func (r *mutationResolver) CreateAccount(c context.Context, in AccountInput) (*Account, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	a, err := r.server.accountClient.PostAccount(c, in.Name, in.Email, in.Password)
	if err != nil {
		return nil, handleError("CreateAccount", err)
	}

	return &Account{
		ID:    a.ID,
		Name:  a.Name,
		Email: a.Email,
	}, nil
}

func (r *mutationResolver) CreateProduct(c context.Context, in ProductInput) (*Product, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(c, in.Name, in.Description, in.Price, in.Image)
	if err != nil {
		return nil, handleError("CreateProduct", err)
	}

	return &Product{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Image:       p.Image,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var products []*productModel.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, &productModel.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}

	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		return nil, handleError("CreateOrder.PostOrder", err)
	}

	var orderProducts []*OrderedProduct
	for _, p := range o.Products {
		prodDetail, err := r.server.catalogClient.GetProduct(ctx, p.ID)
		if err != nil {
			handleError(fmt.Sprintf("GetProduct ID=%s", p.ID), err)
			continue
		}
		orderProducts = append(orderProducts, &OrderedProduct{
			ID:          prodDetail.Id,
			Name:        prodDetail.Name,
			Description: prodDetail.Description,
			Price:       prodDetail.Price,
			Quantity:    int(p.Quantity),
		})
	}

	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
		Products:   orderProducts,
	}, nil
}

func (r *mutationResolver) CreateReview(ctx context.Context, in ReviewInput) (*Review, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if in.Rating < 1 || in.Rating > 5 {
		return nil, ErrInvalidParameter
	}

	review, err := r.server.reviewClient.PostReview(ctx, in.ProductID, in.AccountID, *in.Content, int32(in.Rating))
	if err != nil {
		return nil, handleError("CreateReview", err)
	}

	product, err := r.server.catalogClient.GetProduct(ctx, in.ProductID);
	if err != nil {
		return  nil, handleError("CreateReview", err)
	}

	account, err := r.server.accountClient.GetAccount(ctx, in.AccountID)
	if err != nil {
		return nil, handleError("CreateReview", err)
	}

	return &Review{
		ID:        review.ID,
		Rating:    int(review.Rating),
		Content:   &review.Content,
		CreatedAt: review.CreatedAt,
		Product: &Product{
			ID: product.Id,
			Name: product.Name,
		Price: float64(product.Price),
			Image: product.Image,
		},
		Account: &Account{
			ID: account.ID,
			Name: account.Name,
			Email: account.Email,
		},
	}, nil
}

func (r *mutationResolver) DeleteProduct(c context.Context, id string) (*DeleteResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	result, err := r.server.catalogClient.DeleteProduct(c, id)
	if err != nil {
		handleError("DeleteProduct", err)
		return &DeleteResponse{
			Success:   false,
			Message:   "Failed to delete product",
			DeletedID: id,
		}, nil
	}

	return &DeleteResponse{
		Success:   result.Success,
		Message:   result.Message,
		DeletedID: id,
	}, nil
}

func (r *mutationResolver) Login(c context.Context, in LoginInput) (*AuthResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	token, err := r.server.authClient.Login(c, in.Email, in.Password)
	if err != nil {
		return nil, handleError("Login", err)
	}

	return &AuthResponse{
		ID:    token.Auth.Id,
		Email: token.Auth.Email,
		BackendToken: &Token{
			AccessToken:  token.Auth.Token.AccessToken,
			RefreshToken: token.Auth.Token.RefreshToken,
			ExpiresIn:    int(token.Auth.Token.ExpiresAt),
		},
	}, nil
}

func (r *mutationResolver) RefreshToken(c context.Context, refreshToken string) (*Token, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	newToken, err := r.server.authClient.RefreshToken(c, refreshToken)
	if err != nil {
		return nil, handleError("RefreshToken", err)
	}

	return &Token{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresIn:    int(newToken.ExpiresAt),
	}, nil
}

func (r *mutationResolver) EditProduct(c context.Context, id string, in ProductInput) (*Product, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.EditProduct(c, id, in.Name, in.Description, in.Price, in.Image)
	if err != nil {
		return nil, handleError("EditProduct", err)
	}

	return &Product{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Image:       p.Image,
	}, nil
}

func (r *mutationResolver) DeleteAccount(c context.Context, id string) (*DeleteResponse, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	p, err := r.server.accountClient.DeleteAccount(c, id)
	if err != nil {
		return nil, handleError("DeleteAccount", err)
	}

	return &DeleteResponse{
		DeletedID: p.DeletedID,
		Success:   p.Success,
		Message:   p.Message,
	}, nil
}

func (r *mutationResolver) EditAccount(c context.Context, id string, in EditeAccountInput) (*Account, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	if in.Name == nil {
		empty := ""
		in.Name = &empty
	}
	if in.Email == nil {
		empty := ""
		in.Email = &empty
	}
	if in.Password == nil {
		empty := ""
		in.Password = &empty
	}

	a, err := r.server.accountClient.EditAccount(c, id, *in.Name, *in.Email, *in.Password)
	if err != nil {
		return nil, handleError("EditAccount", err)
	}

	return &Account{
		ID:    a.ID,
		Name:  a.Name,
		Email: a.Email,
	}, nil
}
