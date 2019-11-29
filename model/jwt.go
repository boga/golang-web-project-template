package model

import "github.com/dgrijalva/jwt-go"

type AuthJWT struct {
	AuthIdentityID int `json:"i"`
	jwt.StandardClaims
}

type RefreshJWT struct {
	ID             int `json:"id"`
	AuthIdentityID int `json:"i"`
	jwt.StandardClaims
}
