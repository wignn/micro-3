package repository

import "context"

type Repository interface {
	Close()
	PutOrder(c *context.Context, order )
}