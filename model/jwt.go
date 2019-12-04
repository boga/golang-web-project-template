package model

import (
	"github.com/dgrijalva/jwt-go"
)

type JWTType string

const (
	JWTTypeAuth    JWTType = "a"
	JWTTypeRefresh JWTType = "r"
	JWTTypeTOTP    JWTType = "t"
)

type JWT struct {
	AuthIdentityID int `json:"i"`
	jwt.StandardClaims
	Type JWTType `json:"t"`
}

type AuthJWT struct {
	JWT
}

type RefreshJWT struct {
	ID int `json:"id"`
	JWT
}

type TOTPJWT struct {
	JWT
}
