package resolvers

import (
	"context"
	"fmt"
	"net/http"

	"cipherassets.core/gql"
	"cipherassets.core/gql/errors"
	"cipherassets.core/model"
)

func (r *mutationResolver) TotpVerify(ctx context.Context, data gql.TOTPVerifyInput) (*gql.SigninResponse, error) {
	i := ctx.Value(AuthIdentityContextKey).(*model.AuthIdentity)
	user, err := r.store.UserStore.FindUserByID(i.UserID)
	if err != nil {
		return nil, err
	}

	if !user.TOTPEnabled {
		return nil, errors.NewApiError(nil, "totp-not-enabled", "totp is not enabled for user", http.StatusBadRequest)
	}

	err = r.store.UserStore.ValidateTOTPCode(user, data.Code, false)
	if err != nil {
		return nil, fmt.Errorf("TOTP not valid for user (%d): %w", user.ID, err)
	}
	if err := r.store.UserStore.SaveUser(user); err != nil {
		return nil, fmt.Errorf("can't save TOTP enabled to user (%d): %w", user.ID, err)
	}

	authJWT, err := r.store.JWTStore.MakeAuthJWT(i)
	if err != nil {
		return nil, fmt.Errorf("can't make auth jwt for i %d: %w", i.ID, err)
	}
	authJWTStr, err := r.store.JWTStore.GetJWTString(authJWT)
	if err != nil {
		return nil, fmt.Errorf("can't make auth jwt string for i %d: %w", i.ID, err)
	}

	refreshJWT, err := r.store.JWTStore.MakeRefreshJWT(i)
	if err != nil {
		return nil, fmt.Errorf("can't make refresh jwt for i %d: %w", i.ID, err)
	}
	refreshJWTStr, err := r.store.JWTStore.GetJWTString(refreshJWT)
	if err != nil {
		return nil, fmt.Errorf("can't make refresh jwt string for i %d: %w", i.ID, err)
	}

	return &gql.SigninResponse{
		AuthToken:    *authJWTStr,
		RefreshToken: *refreshJWTStr,
	}, nil
}
