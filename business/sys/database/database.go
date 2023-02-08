// Package database provides support for access the database.
package database

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"net/url"
	"time"
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

func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	for i := 1; ; i++ {
		if err := db.Ping(); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ticker.Reset(time.Duration(i) * 100 * time.Millisecond)
		}
	}
}
