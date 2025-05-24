package main

import "context"

type accountResolver struct {
	server *GraphQLServer
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	
}