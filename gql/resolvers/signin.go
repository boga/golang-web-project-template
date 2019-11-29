package resolvers

import (
	"context"
	"fmt"

	"cipherassets.core/gql"
)

func (r *mutationResolver) Signin(ctx context.Context, creds gql.SigninInput) (*gql.SigninResponse, error) {
	i, err := r.store.UserStore.FindIdentityByUID(creds.Email)
	if err != nil {
		return nil, fmt.Errorf("identity not found: %w", err)
	}

	// TODO: add password verification

	authJWT, err := r.store.JWTStore.MakeAuthJWT(i)
	if err != nil {
		return nil, fmt.Errorf("can't make auth jwt for identity %d: %w", i.ID, err)
	}
	authJWTStr, err := r.store.JWTStore.GetJWTString(authJWT)
	if err != nil {
		return nil, fmt.Errorf("can't make auth jwt string for identity %d: %w", i.ID, err)
	}

	refreshJWT, err := r.store.JWTStore.MakeRefreshJWT(i)
	if err != nil {
		return nil, fmt.Errorf("can't make refresh jwt for identity %d: %w", i.ID, err)
	}
	refreshJWTStr, err := r.store.JWTStore.GetJWTString(refreshJWT)
	if err != nil {
		return nil, fmt.Errorf("can't make refresh jwt string for identity %d: %w", i.ID, err)
	}

	return &gql.SigninResponse{
		AuthToken:    *authJWTStr,
		RefreshToken: *refreshJWTStr,
	}, nil
}
