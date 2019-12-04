package resolvers

import (
	"context"
	"fmt"
	"net/http"

	"cipherassets.core/gql"
	"cipherassets.core/gql/errors"
)

func (r *mutationResolver) Signin(ctx context.Context, creds gql.SigninInput) (*gql.SigninResponse, error) {
	i, err := r.store.UserStore.FindIdentityByUID(creds.Email)
	if err != nil {
		return nil, fmt.Errorf("identity not found: %w", err)
	}

	if err := r.store.UserStore.CheckPassword(i, creds.Password); err != nil {
		return nil, errors.NewApiError(err, "wrong-password", "wrong password", http.StatusBadRequest)
	}

	u, err := r.store.UserStore.FindUserByID(i.UserID)
	if err != nil {
		return nil, err
	}
	if u.TOTPEnabled {
		totpJWT, err := r.store.JWTStore.MakeTOTPJWT(i)
		if err != nil {
			return nil, fmt.Errorf("can't make totp jwt for identity %d: %w", i.ID, err)
		}
		totpJWTStr, err := r.store.JWTStore.GetJWTString(totpJWT)
		if err != nil {
			return nil, fmt.Errorf("can't make totp jwt string for identity %d: %w", i.ID, err)
		}
		return &gql.SigninResponse{
			TotpToken: totpJWTStr,
		}, nil
	} else {
		refreshJWT, err := r.store.JWTStore.MakeRefreshJWT(i)
		if err != nil {
			return nil, fmt.Errorf("can't make refresh jwt for identity %d: %w", i.ID, err)
		}
		refreshJWTStr, err := r.store.JWTStore.GetJWTString(refreshJWT)
		if err != nil {
			return nil, fmt.Errorf("can't make refresh jwt string for identity %d: %w", i.ID, err)
		}

		return &gql.SigninResponse{
			RefreshToken: refreshJWTStr,
		}, nil
	}

}
