package userdb

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Repository struct {
	db  *sqlx.DB
	log *zap.Logger
}

const (
	createQuery  = `INSERT INTO users(name, email, password, role) VALUES ($1, $2, $3, $4)`
	deleteQuery  = `DELETE FROM users WHERE id = $1`
	updateQuery  = `UPDATE users SET name=?, email=?, role=?, password=? WHERE id = ?`
	getByIDQuery = `SELECT name, email, password, role FROM users WHERE id = ?`
)

func NewRepository(db *sqlx.DB, log *zap.Logger) Repository {
	return Repository{db: db, log: log}
}

func (r Repository) Create(ctx context.Context, nu NewUser) (User, error) {
	var user User
	if err := r.db.QueryRowxContext(ctx, createQuery, nu.Name, nu.Email, nu.Password, nu.Role).StructScan(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r Repository) Delete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("something weird happened")
	}
	return nil
}
