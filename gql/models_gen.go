// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"cipherassets.core/model"
)

type SigninInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	User *model.User `json:"user"`
}
