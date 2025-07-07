package repository

import (
	"context"
	"database/sql"
	"github.com/wignn/micro-3/cart/model"
	_ "github.com/lib/pq"
)

type CartRepository interface {
	Close()
	PutCart(ctx context.Context, o *model.CartPutRequest) error
	GetCartByAccount(ctx context.Context, accountID string) (*model.Cart, error)
	DeleteCart(ctx context.Context, id string) error
	GetCartByID(ctx context.Context, id string) (*model.Cart, error)
	UpdateCartItem(ctx context.Context, cartID, productID string, quantity int) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (CartRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close() {
	_ = r.db.Close()
}

func (r *PostgresRepository) PutCart(ctx context.Context, o *model.CartPutRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var cartID string

	err = tx.QueryRowContext(ctx, `
		SELECT id FROM carts WHERE account_id = $1
	`, o.AccountID).Scan(&cartID)

	if err == sql.ErrNoRows {
		cartID = o.ID
		_, err = tx.ExecContext(ctx, `
			INSERT INTO carts (id, account_id, created_at)
			VALUES ($1, $2, $3)
		`, cartID, o.AccountID, o.CreatedAt)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`, cartID, o.ProductID, o.Quantity)
	if err != nil {
		return err
	}

	return tx.Commit()
}


func (r *PostgresRepository) GetCartByID(ctx context.Context, id string) (*model.Cart, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, account_id, created_at
		FROM carts
		WHERE id = $1
	`, id)

	var cart model.Cart
	err := row.Scan(&cart.ID, &cart.AccountID, &cart.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT product_id, quantity
		FROM cart_items
		WHERE cart_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart.Products = []model.CartProduct{}
	for rows.Next() {
		var product model.CartProduct
		if err := rows.Scan(&product.ProductID, &product.Quantity); err != nil {
			return nil, err
		}
		cart.Products = append(cart.Products, product)
	}

	return &cart, nil
}

func (r *PostgresRepository) GetCartByAccount(ctx context.Context, accountID string) (*model.Cart, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, account_id, created_at
		FROM carts
		WHERE account_id = $1
	`, accountID)

	var cart model.Cart
	err := row.Scan(&cart.ID, &cart.AccountID, &cart.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT product_id, quantity
		FROM cart_items
		WHERE cart_id = $1
	`, cart.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart.Products = []model.CartProduct{}
	for rows.Next() {
		var product model.CartProduct
		if err := rows.Scan(&product.ProductID, &product.Quantity); err != nil {
			return nil, err
		}
		cart.Products = append(cart.Products, product)
	}

	return &cart, nil
}

func (r *PostgresRepository) DeleteCart(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM carts WHERE id = $1
	`, id)
	return err
}

func (r *PostgresRepository) UpdateCartItem(ctx context.Context, cartID, productID string, quantity int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE cart_items
		SET quantity = $1
		WHERE cart_id = $2 AND product_id = $3
	`, quantity, cartID, productID)
	return err
}
