package cipherassets_core

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DB *sqlx.DB
}

func NewDB(c *Config) (*DB, error) {
	// db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s", c.Database.driver,
		c.Database.user, c.Database.pass, c.Database.host, c.Database.port, c.Database.name)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DB: %w", err)
	}

	return &DB{DB: db}, nil
}
