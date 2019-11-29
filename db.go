package cipherassets_core

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDB(c *Config) (*DB, error) {
	// db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	dbSQLX, err := sqlx.Connect(c.Database.driver, fmt.Sprintf("%s:%s@(%s:%d)/%s",
		c.Database.user, c.Database.pass, c.Database.host, c.Database.port, c.Database.name))

	if err != nil {
		return nil, fmt.Errorf("can't connect to DB: %w", err)
	}
	db := DB{dbSQLX}

	return &db, nil
}

func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	err := db.DB.Get(dest, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NotFoundError{err: err}
		}

		return err
	}

	return nil
}

type NotFoundError struct {
	err error
}

func (e NotFoundError) Error() string {
	return "row(s) not found"
}

func (e NotFoundError) Unwrap() error {
	return e.err
}
