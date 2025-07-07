package client

import (
	"context"
	"log"

	"github.com/wignn/micro-3/catalog/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CatalogClient struct {
	conn    *grpc.ClientConn
	service genproto.CatalogServiceClient
}

func NewClient(url string) (*CatalogClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to catalog service: %v\n", err)
		return nil, err
	}

	c := genproto.NewCatalogServiceClient(conn)
	return &CatalogClient{conn: conn, service: c}, nil
}

func (cl *CatalogClient) Close() {
	cl.conn.Close()
}

func (cl *CatalogClient) PostProduct(c context.Context, name, description string, price float64, image string) (*genproto.Product, error) {

	r, err := cl.service.PostProduct(
		c,
		&genproto.PostProductRequest{Name: name, Description: description, Price: price, Image: image},
	)

	if err != nil {
		log.Printf("failed to post product: %v\n", err)
		return nil, err
	}

	return r.Product, nil
}

func (cl *CatalogClient) GetProduct(c context.Context, id string) (*genproto.Product, error) {
	r, err := cl.service.GetProduct(
		c,
		&genproto.GetProductRequest{Id: id},
	)
	if err != nil {
		log.Printf("failed to get product: %v\n", err)
		return nil, err
	}

	return r.Product, nil
}

func (cl *CatalogClient) GetProducts(c context.Context, skip, take uint64, ids []string, query string) ([]*genproto.Product, error) {
	r, err := cl.service.GetProducts(
		c,
		&genproto.GetProductsRequest{Skip: skip, Take: take, Ids: ids, Query: query},
	)
	if err != nil {
		log.Printf("failed to get products: %v\n", err)
		return nil, err
	}

	var products []*genproto.Product

	for _, p := range r.Products {
		products = append(products, &genproto.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Image:       p.Image,
		})
	}

	return products, nil
}

func (cl *CatalogClient) DeleteProduct(c context.Context, id string) (*genproto.DeleteProductResponse, error) {
	result, err := cl.service.DeleteProduct(
		c,
		&genproto.DeleteProductRequest{Id: id},
	)
	if err != nil {
		log.Printf("failed to delete product: %v\n", err)
		return nil, err
	}
	return &genproto.DeleteProductResponse{
		DeletedID: id,
		Success:   result.Success,
		Message:   result.Message,
	}, nil
}

func (cl *CatalogClient) EditProduct(c context.Context, id string, name, description string, price float64, image string) (*genproto.Product, error) {
	r, err := cl.service.EditProduct(
		c,
		&genproto.EditProductRequest{
			Id:          id,
			Name:        name,
			Description: description,
			Price:       price,
			Image:       image,
		},
	)
	if err != nil {
		log.Printf("failed to edit product: %v\n", err)
		return nil, err
	}

	return r.Product, nil
}
