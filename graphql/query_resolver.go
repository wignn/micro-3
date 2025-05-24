package main

import "context"

type queryResolver struct {
	server *GraphQLServer
}

func (r *queryResolver) Account(c context.Context,pagination *PaginationInput, id *string) ([]*Account, error) {
	
}
package main

import "context"

type queryResolver struct {
	server *GraphQLServer
}

// func (r *queryResolver) Account(c context.Context,pagination *PaginationInput, id *string) ([]*Account, error) {

// }
// func (r *queryResolver) Product(c context.Context,pagination *PaginationInput, query *string,id *string) ([]*Product, error) {

// }
// func (r *queryResolver) Account(c context.Context,pagination *PaginationInput, query *string,id *string) ([]*Order, error) {

// }