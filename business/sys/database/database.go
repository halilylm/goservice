// Package database provides support for access the database.
package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"net/url"
)

type Config struct {
	User         string
	Password     string
	Name         string
	Host         string
	MaxOpenConns int
	MaxIdleConns int
	DisableTLS   bool
}

func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "true"
	if cfg.DisableTLS {
		sslMode = "false"
	}
	q := make(url.Values)
	q.Set("tls", sslMode)
	dsn := url.URL{
		Scheme:   "mysql",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawPath:  cfg.Name,
		RawQuery: q.Encode(),
	}
	conn, err := sqlx.Connect("mysql", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("connecting: %w", err)
	}
	return conn, nil
}
