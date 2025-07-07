package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/wignn/micro-3/order/model"
)

type OrderRepository interface {
	Close()
	PutOrder(c context.Context, o *model.Order) error
	GetOrdersForAccount(c context.Context, accountID string) ([]*model.Order, error)
	DeleteOrder(c context.Context, id string) error
}

type postgresRepository struct {
	db *sql.DB
}


func NewOrderPostgresRepository(url string) (OrderRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(c context.Context, o *model.Order) (err error) {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// order creation
	_, err = tx.ExecContext(
		c,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)
	if err != nil {
		return
	}

	// order products insertion
	stmt, _ := tx.PrepareContext(c, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	for _, p := range o.Products {
		_, err = stmt.ExecContext(c, o.ID, p.ID, p.Quantity)
		if err != nil {
			return
		}
	}
	_, err = stmt.ExecContext(c)
	if err != nil {
		return
	}
	stmt.Close()

	return 
}



func (r *postgresRepository) GetOrdersForAccount(c context.Context, accountID string) ([]*model.Order, error) {
	rows, err := r.db.QueryContext(
		c,
		`SELECT
      o.id,
      o.created_at,
      o.account_id,
      o.total_price::money::numeric::float8,
      op.product_id,
      op.quantity
    FROM orders o JOIN order_products op ON (o.id = op.order_id)
    WHERE o.account_id = $1
    ORDER BY o.id`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*model.Order{}
	order := &model.Order{}
	lastOrder := &model.Order{}
	orderedProduct := &model.OrderedProduct{}
	products := []model.OrderedProduct{}

	// Scan rows into Order structs
	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}
		// Scan order
		if lastOrder.ID != "" && lastOrder.ID != order.ID {
			newOrder := &model.Order{
				ID:         lastOrder.ID,
				AccountID:  lastOrder.AccountID,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   products,
			}
			orders = append(orders, newOrder)
			products = []model.OrderedProduct{}
		}
		// Scan products
		products = append(products, model.OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})

		*lastOrder = *order
	}

	// Add last order (or first :D)
	if lastOrder != nil {
		newOrder := &model.Order{
			ID:         lastOrder.ID,
			AccountID:  lastOrder.AccountID,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   products,
		}
		orders = append(orders, newOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}


func (r *postgresRepository) DeleteOrder(c context.Context, id string) error {
	res, err := r.db.ExecContext(
		c,
		"DELETE FROM orders WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no order found with id %s", id)
	}

	return nil
}