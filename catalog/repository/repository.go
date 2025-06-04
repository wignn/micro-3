package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/wignn/micro-3/catalog/model"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type CatalogRepository interface {
	Close()
	PutProduct(c context.Context, p *model.Product) error
	GetProductByID(c context.Context, id string) (*model.Product, error)
	ListProducts(c context.Context, skip uint64, take uint64) ([]*model.Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]*model.Product, error)
	SearchProducts(c context.Context, query string, skip uint64, take uint64) ([]*model.Product, error)
	DeletedProduct(c context.Context, id string) error
}

type elasticRepository struct {
	client *elastic.Client
}



func NewElasticRepository(url string) (CatalogRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {}

func (r *elasticRepository) PutProduct(ctx context.Context, p *model.Product) error {
	_, err := r.client.Index().
		Index("catalog").
		Type("product").
		Id(p.ID).
		BodyJson(model.ProductDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Image:       p.Image,
		}).
		Do(ctx)
	return err
}

func (r *elasticRepository) GetProductByID(c context.Context, id string) (*model.Product, error) {
	res, err := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(c)
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}
	p := model.ProductDocument{}
	if err = json.Unmarshal(*res.Source, &p); err != nil {
		return nil, err
	}
	return &model.Product{
		ID:          res.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Image:       p.Image,
	}, nil
}

func (r *elasticRepository) ListProducts(c context.Context, skip, take uint64) ([]*model.Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).Size(int(take)).
		Do(c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []*model.Product
	for _, hit := range res.Hits.Hits {
		p := model.ProductDocument{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, &model.Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Image:       p.Image,
			})
		}
	}
	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(c context.Context, ids []string) ([]*model.Product, error) {
	var items []*elastic.MultiGetItem
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().
			Index("catalog").
			Type("product").
			Id(id))
	}
	res, err := r.client.MultiGet().
		Add(items...).
		Do(c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []*model.Product
	for _, doc := range res.Docs {
		if doc.Found {
			p := model.ProductDocument{}
			if err = json.Unmarshal(*doc.Source, &p); err == nil {
				products = append(products, &model.Product{
					ID:          doc.Id,
					Name:        p.Name,
					Description: p.Description,
					Price:       p.Price,
					Image:       p.Image,
				})
			}
		}
	}
	return products, nil
}

func (r *elasticRepository) SearchProducts(c context.Context, query string, skip, take uint64) ([]*model.Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).Size(int(take)).
		Do(c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []*model.Product
	for _, hit := range res.Hits.Hits {
		p := model.ProductDocument{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, &model.Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Image:       p.Image,
			})
		}
	}
	return products, nil
}

func (r *elasticRepository) DeletedProduct(c context.Context, id string) error {
	_, err := r.client.Delete().
		Index("catalog").
		Type("product").
		Id(id).
		Do(c)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}