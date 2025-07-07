package main

import (
	"context"
	"log"
	"time"
)

type accountResolver struct {
	server *GraphQLServer
}

func (r *accountResolver) Orders(c context.Context, o *Account) ([]*Order, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrdersForAccount(c, o.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*Order
	for _, o := range orderList {
		var products []*OrderedProduct
		for _, p := range o.Products {
			products = append(products, &OrderedProduct{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID:         o.ID,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
			Products:   products,
			
		})
	}

	return orders, nil
}
