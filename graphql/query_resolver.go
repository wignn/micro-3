package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *GraphQLServer
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get single
	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{{
			ID:    r.ID,
			Name:  r.Name,
			Email: r.Email,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accounts []*Account
	for _, a := range accountList {
		account := &Account{
			ID:    a.ID,
			Name:  a.Name,
			Email: a.Email,
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *queryResolver) Products(c context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	// Get single
	if id != nil {
		r, err := r.server.catalogClient.GetProduct(c, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{{
			ID:          r.Id,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
			Image:       r.Image,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}
	
	productList, err := r.server.catalogClient.GetProducts(c, skip, take, nil, q)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*Product

	for _, a := range productList {
		products = append(products,
			&Product{
				ID:          a.Id,
				Name:        a.Name,
				Description: a.Description,
				Price:       a.Price,
				Image:       a.Image,
			},
		)
	}

	return products, nil
}
func (r *queryResolver) Reviews(c context.Context, pagination *PaginationInput, id *string) ([]*Review, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	if pagination == nil {
		rv, err := r.server.reviewClient.GetReview(c, *id)
		if err != nil {
			log.Println("GetReview error:", err)
			return nil, err
		}

		var createdAt time.Time

		if err := createdAt.UnmarshalBinary(rv.CreatedAt); err != nil {
			log.Println("UnmarshalBinary error on review ID", rv.Id, ":", err)
			return nil, err
		} 

		return []*Review{{
			ID:        rv.Id,
			ProductID: rv.ProductId,
			AccountID: rv.AccountId,
			Rating:    int(rv.Rating),
			Content:   &rv.Content,
			CreatedAt: createdAt,
		}}, nil
	}

	skip, take := pagination.bounds()
	reviewList, err := r.server.reviewClient.GetReviews(c, *id, skip, take)
	if err != nil {
		log.Println("GetReviews error:", err)
		return nil, err
	}

	var reviews []*Review
	for _, a := range reviewList {
		var createdAt time.Time
		if err := createdAt.UnmarshalBinary(a.CreatedAt); err != nil {
			log.Println("UnmarshalBinary error on review ID", a.Id, ":", err)
			continue
		}

		reviews = append(reviews, &Review{
			ID:        a.Id,
			ProductID: a.ProductId,
			AccountID: a.AccountId,
			Rating:    int(a.Rating),
			Content:   &a.Content,
			CreatedAt: createdAt,
		})
	}

	return reviews, nil
}


func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}
	return skipValue, takeValue
}
