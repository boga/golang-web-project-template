package store

import (
	"fmt"

	core "cipherassets.core"
	"cipherassets.core/model"
)

type UserStore struct {
	db *core.DB
}

func (s *UserStore) GetUsers() ([]model.User, error) {
	var people []model.User
	if err := s.db.DB.Select(&people, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("can't get users: %w", err)
	}

	return people, nil
}
