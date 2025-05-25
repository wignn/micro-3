package service

import (
	"context"
	"github.com/wignn/micro-3/catalog/model"
	"github.com/wignn/micro-3/catalog/repository"
)

type CatalogService interface {
	PostProduct(c context.Context, name, description string, price float64) (*model.Product, error)
	GetProduct(c context.Context, id string) (*model.Product, error)
	GetProducts(c context.Context, skip uint64, take uint64) ([]*model.Product, error)
	GetProductsByIDs(c context.Context, ids []string) ([]*model.Product, error)
	SearchProducts(c context.Context, query string, skip uint64, take uint64) ([]*model.Product, error)
}

type catalogService struct {
	repository repository.CatalogRepository
}

func NewCatalogService(r repository.CatalogRepository) CatalogService {
	return &catalogService{r}
}

func (s *catalogService) PostProduct(c context.Context, name, description string, price float64) (*model.Product, error) {
	p := &model.Product{
		Name:        name,
		Description: description,
		Price:       price,
	}

	if err := s.repository.PutProduct(c, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *catalogService) GetProduct(c context.Context, id string) (*model.Product, error) {
	return s.repository.GetProductByID(c, id)
}

func (s *catalogService) GetProducts(c context.Context, skip uint64, take uint64) ([]*model.Product, error) {
	if skip > 100 || (skip == 100 && take > 100) {
		take = 100
	}
	return s.repository.ListProducts(c, skip, take)
}

func (s *catalogService) GetProductsByIDs(c context.Context, ids []string) ([]*model.Product, error) {
	return s.repository.ListProductsWithIDs(c, ids)
}

func (s *catalogService) SearchProducts(c context.Context, query string, skip uint64, take uint64) ([]*model.Product, error) {
	if skip > 100 || (skip == 100 && take > 100) {
		take = 100
	}
	return s.repository.SearchProducts(c, query, skip, take)
}