package store

import (
	"errors"
	"fmt"

	core "cipherassets.core"
	"cipherassets.core/model"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *core.DB
}

const PasswordDefaultCost = 10

func (s *UserStore) NewIdentity(uid string, password string, userID int) (*model.AuthIdentity, error) {
	i := &model.AuthIdentity{
		UID:    uid,
		UserID: userID,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password) /* bcrypt.DefaultCost */, PasswordDefaultCost)
	if err != nil {
		return nil, fmt.Errorf("can't hash password: %w", err)
	}
	hashStr := string(hash)
	i.Password = &hashStr

	return i, nil
}

func (s *UserStore) NewUser() (*model.User, error) {
	u := &model.User{}

	return u, nil
}

func (s *UserStore) SaveIdentity(i *model.AuthIdentity) error {
	var query string
	var err error
	if i.ID == 0 {
		query = "INSERT INTO auth_identities (user_id, uid, password) VALUES (:user_id, :uid, :password)"
	} else {
		query = "UPDATE auth_identities SET user_id = ?, uid = ?, password = ? WHERE id = :id"
	}
	result, err := s.db.NamedExec(query, i)
	if err != nil {
		return fmt.Errorf("can't save identity to db: %w. Identity: %+v", err, i)
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("can't get ID of inserted identity: %w. Identity: %+v", err, i)
	}

	if i.ID == 0 {
		i.ID = int(lastInsertId)
	}

	return nil
}

func (s *UserStore) SaveUser(u *model.User) error {
	var query string
	var err error
	if u.ID == 0 {
		query = "INSERT INTO users (name) VALUES (:name)"
	} else {
		query = "UPDATE users SET name = ? WHERE id = :id"
	}
	result, err := s.db.NamedExec(query, u)
	if err != nil {
		return fmt.Errorf("can't save user to db: %w. User: %+v", err, u)
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("can't get ID of inserted user: %w. User: %+v", err, u)
	}

	if u.ID == 0 {
		u.ID = int(lastInsertId)
	}

	return nil
}

func (s *UserStore) FindIdentityByUID(uid string) (*model.AuthIdentity, error) {
	i := model.AuthIdentity{}
	if err := s.db.Get(&i, "SELECT * FROM auth_identities WHERE uid = ?", uid); err != nil {
		var nfe core.NotFoundError
		if errors.As(err, &nfe) {
			return nil, nil
		}

		return nil, fmt.Errorf("can't get identity by uid='%s': %w", uid, err)
	}

	return &i, nil
}

func (s *UserStore) GetUsers() ([]model.User, error) {
	var people []model.User
	if err := s.db.Select(&people, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("can't get users: %w", err)
	}

	return people, nil
}
