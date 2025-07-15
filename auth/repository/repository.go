package repository

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/wignn/micro-3/auth/model"
)

type AuthRepository interface {
	Close()
	GetAccount(c context.Context, email string) (*model.AuthResponseRepository, error)
}

type authRepository struct {
	db *sql.DB
}


func NewAuthPostgresRepository(url string) (AuthRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &authRepository{db: db}, nil
}

func (r *authRepository) Close() {
	r.db.Close()
}

func (r *authRepository) GetAccount(c context.Context, email string) (*model.AuthResponseRepository, error) {
	var account model.AuthResponseRepository
	err := r.db.QueryRowContext(c, "SELECT id, email, password FROM accounts WHERE email = $1", email).Scan(&account.ID, &account.Email, &account.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}
	return &account, nil
}