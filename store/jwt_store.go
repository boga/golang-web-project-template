package store

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	core "cipherassets.core"
	"cipherassets.core/model"
)

type JWTStore struct {
	config *core.Config
	db     *core.DB
}

func (s *JWTStore) GetJWTString(t jwt.Claims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, t)
	var jwtSecret = []byte(s.config.JWT.Secret)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *JWTStore) ParseJWTString(tokenStr *string, claims jwt.Claims) error {
	var err error
	tkn := &jwt.Token{
		Claims: claims,
	}
	tkn, err = jwt.ParseWithClaims(*tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})
	_ = err
	if tkn == nil || !tkn.Valid {
		return jwt.ErrInvalidKey
	}

	return nil
}

func (s *JWTStore) MakeAuthJWT(identity *model.AuthIdentity) (*model.AuthJWT, error) {
	t := &model.AuthJWT{}
	t.AuthIdentityID = identity.ID
	t.ExpiresAt = time.Now().Add(*s.config.JWT.AuthTokenTTL).Unix()
	t.Type = model.JWTTypeAuth

	return t, nil
}

func (s *JWTStore) MakeRefreshJWT(identity *model.AuthIdentity) (*model.RefreshJWT, error) {
	t := &model.RefreshJWT{}
	t.AuthIdentityID = identity.ID
	t.ExpiresAt = time.Now().Add(time.Hour * 24 * 365).Unix()
	t.Type = model.JWTTypeRefresh
	insertResult := s.db.MustExec("INSERT INTO refresh_tokens (auth_identity_id) VALUES (?)", identity.ID)
	lastInsertId, err := insertResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("can't insert refresh token: %w", err)
	}
	t.ID = int(lastInsertId)

	return t, nil
}

func (s *JWTStore) MakeTOTPJWT(identity *model.AuthIdentity) (*model.TOTPJWT, error) {
	t := &model.TOTPJWT{}
	t.AuthIdentityID = identity.ID
	t.ExpiresAt = time.Now().Add(time.Minute * 15).Unix()
	t.Type = model.JWTTypeTOTP

	return t, nil
}
