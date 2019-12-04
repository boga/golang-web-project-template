package store

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	core "cipherassets.core"
	"cipherassets.core/model"
)

type UserStore struct {
	db *core.DB
}

const PasswordDefaultCost = 10
const backupCodeLen = 10

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

func (s *UserStore) ValidateTOTPCode(user *model.User, code string, noBackupCodesIsOK bool) error {
	ok := false

	if len(code) == backupCodeLen {
		backupCodes := strings.Split(*user.TOTPBackupCodes, "\n")
		for i := 0; !ok && i < len(backupCodes); i++ {
			ok = code == backupCodes[i]
		}
	} else {
		ok = totp.Validate(code, *user.TOTPSecret)
	}

	if !ok {
		return errors.New("code is wrong")
	}

	if err := s.RemoveBackupCode(user, code); err != nil {
		if _, ok = err.(NoBackupCodesError); ok && noBackupCodesIsOK {
			return nil
		}
		return err
	}

	return nil
}

func (s *UserStore) CheckPassword(identity *model.AuthIdentity, password string) error {
	if identity.Password == nil {
		return fmt.Errorf("password no set")
	}
	hash := []byte(*identity.Password)
	pass := []byte(password)
	if err := bcrypt.CompareHashAndPassword(hash, pass); err != nil {
		return fmt.Errorf("password not valid: %w", err)
	}

	return nil
}

func (s *UserStore) RemoveBackupCode(user *model.User, code string) error {
	if *user.TOTPBackupCodes == "" {
		return NoBackupCodesError{}
	}

	backupCodes := strings.Split(*user.TOTPBackupCodes, "\n")
	idx := -1
	for i := 0; i < len(backupCodes); i++ {
		if code == backupCodes[i] && len(backupCodes[i]) != 0 {
			idx = i
		}
	}
	if idx > -1 {
		backupCodes[idx] = backupCodes[len(backupCodes)-1]
		*user.TOTPBackupCodes = strings.Join(backupCodes[0:len(backupCodes)-1], "\n")
	}

	return nil
}

type NoBackupCodesError struct {
}

func (e NoBackupCodesError) Error() string {
	return "No backup codes"
}

func (s *UserStore) SaveUser(u *model.User) error {
	var query string
	var err error
	if u.ID == 0 {
		query = `
			INSERT INTO users (name, totp_backup_codes, totp_enabled, totp_secret) 
			           VALUES (:name, :totp_backup_codes, :totp_enabled, :totp_secret)
	   `
	} else {
		query = `
			UPDATE users 
			SET name = :name, 
			 	totp_backup_codes = :totp_backup_codes,
			 	totp_enabled = :totp_enabled,
			 	totp_secret = :totp_secret
			WHERE id = :id
		`
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

func (s *UserStore) FindIdentityByID(id int) (*model.AuthIdentity, error) {
	i := model.AuthIdentity{}
	if err := s.db.Get(&i, "SELECT * FROM auth_identities WHERE id = ?", id); err != nil {
		var nfe core.NotFoundError
		if errors.As(err, &nfe) {
			return nil, nil
		}

		return nil, fmt.Errorf("can't get identity by id='%d': %w", id, err)
	}

	return &i, nil
}

func (s *UserStore) FindUserByID(id int) (*model.User, error) {
	u := model.User{}
	if err := s.db.Get(&u, "SELECT * FROM users WHERE id = ?", id); err != nil {
		var nfe core.NotFoundError
		if errors.As(err, &nfe) {
			return nil, nil
		}

		return nil, fmt.Errorf("can't get user by id='%d': %w", id, err)
	}

	return &u, nil
}

func (s *UserStore) GetUsers() ([]model.User, error) {
	var people []model.User
	if err := s.db.Select(&people, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("can't get users: %w", err)
	}

	return people, nil
}

func (s *UserStore) NewBackupCodes() []string {
	var backupCodes []string
	letters := []rune("ABCDEFGHJKLMNPQRSTWXYZabcdefghjkmnpqrstwxyz23456789")
	rand.Seed(time.Now().UnixNano())
	for len(backupCodes) < 10 {
		code := ""
		for len(code) < backupCodeLen {
			code += string(letters[rand.Intn(len(letters))])
		}
		backupCodes = append(backupCodes, code)
	}

	return backupCodes
}
