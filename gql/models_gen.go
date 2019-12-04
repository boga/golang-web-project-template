// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"cipherassets.core/model"
)

type GetAuthJWTResponse struct {
	AuthToken string `json:"auth_token"`
}

type SigninInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninResponse struct {
	RefreshToken *string `json:"refresh_token"`
	TotpToken    *string `json:"totp_token"`
}

type SignupInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	User *model.User `json:"user"`
}

type TOTPDisableInput struct {
	Code string `json:"code"`
}

type TOTPDisableResponse struct {
	Success bool `json:"success"`
}

type TOTPGenerateResponse struct {
	Qrcode string `json:"qrcode"`
}

type TOTPSetupInput struct {
	Code string `json:"code"`
}

type TOTPSetupResponse struct {
	BackupCodes []string `json:"backupCodes"`
}

type TOTPVerifyInput struct {
	Code string `json:"code"`
}
