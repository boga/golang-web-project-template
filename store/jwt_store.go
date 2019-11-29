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

func (s *JWTStore) MakeAuthJWT(identity *model.AuthIdentity) (*model.AuthJWT, error) {
	t := &model.AuthJWT{
		AuthIdentityID: identity.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(*s.config.JWT.AuthTokenTTL).Unix(),
		},
	}

	return t, nil
}

func (s *JWTStore) MakeRefreshJWT(identity *model.AuthIdentity) (*model.RefreshJWT, error) {
	t := &model.RefreshJWT{
		AuthIdentityID: identity.ID,
	}
	insertResult := s.db.MustExec("INSERT INTO refresh_tokens (auth_identity_id) VALUES (?)", identity.ID)
	lastInsertId, err := insertResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("can't insert refresh token: %w", err)
	}
	t.ID = int(lastInsertId)

	return t, nil
}
