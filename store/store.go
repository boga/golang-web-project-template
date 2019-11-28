package store

import (
	"fmt"

	core "cipherassets.core"
)

type Store struct {
	UserStore *UserStore
}

func NewStore(config *core.Config) (*Store, error) {
	db, err := core.NewDB(config)
	if err != nil {
		return nil, fmt.Errorf("can't create DB: %w", err)
	}
	return &Store{
		UserStore: &UserStore{db: db},
	}, nil
}
