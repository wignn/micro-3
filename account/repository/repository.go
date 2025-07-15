package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wignn/micro-3/account/model"
)

var (
	ErrNotFound = errors.New("entity not found")
)


type AccountRepository interface {
	Close()	
	PutAccount(c context.Context, a *model.Account) error
	GetAccountById(c context.Context, id string) (*model.Account, error)
	ListAccount(c context.Context, skip uint64, take uint64) ([]*model.Account, error)
	EditAccount(c context.Context, a *model.Account) (*model.Account, error)
	DeleteAccount(c context.Context, id string) error
}


type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Close() {
	if err := r.db.Close(); err != nil {
		panic(err)
	}
}

func (r *PostgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *PostgresRepository) PutAccount(c context.Context, a *model.Account) error {
	_, err := r.db.ExecContext(c, "INSERT INTO accounts (id, name, email, password) VALUES ($1, $2, $3, $4)", a.ID, a.Name, a.Email, a.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetAccountById(c context.Context, id string) (*model.Account, error) {
	row := r.db.QueryRowContext(c, "SELECT id, name, email FROM accounts WHERE id = $1", id)
	a := &model.Account{}
	if err := row.Scan(&a.ID, &a.Name, &a.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return a, nil
}

func (r *PostgresRepository) ListAccount(c context.Context, skip uint64, take uint64) ([]*model.Account, error) {
	rows, err := r.db.QueryContext(c, "SELECT id, name, email FROM accounts OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*model.Account
	
	for rows.Next() {
		a := &model.Account{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Email); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *PostgresRepository) DeleteAccount(c context.Context, id string) error {
	res, err := r.db.ExecContext(c, "DELETE FROM accounts WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no account found with id %s", id)
	}

	return nil
}


func (r *PostgresRepository) EditAccount(c context.Context, a *model.Account) (*model.Account, error) {
	// Ambil data lama dulu
	var old model.Account
	err := r.db.QueryRowContext(c, "SELECT id, name, email, password FROM accounts WHERE id = $1", a.ID).
		Scan(&old.ID, &old.Name, &old.Email, &old.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Ganti field kosong dengan data lama
	if a.Name == "" {
		a.Name = old.Name
	}
	if a.Email == "" {
		a.Email = old.Email
	}
	if a.Password == "" {
		a.Password = old.Password
	}

	// Lakukan update
	_, err = r.db.ExecContext(c,
		"UPDATE accounts SET name = $1, email = $2, password = $3 WHERE id = $4",
		a.Name, a.Email, a.Password, a.ID)
	if err != nil {
		return nil, err
	}

	// Ambil kembali data terbaru
	var updated model.Account
	err = r.db.QueryRowContext(c, "SELECT id, name, email FROM accounts WHERE id = $1", a.ID).
		Scan(&updated.ID, &updated.Name, &updated.Email)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}
