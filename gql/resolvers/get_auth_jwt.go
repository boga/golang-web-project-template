package resolvers

import (
	"context"
	"fmt"

	"cipherassets.core/gql"
	"cipherassets.core/gql/errors"
	"cipherassets.core/model"
)

func (r *queryResolver) GetAuthJwt(ctx context.Context) (*gql.GetAuthJWTResponse, error) {
	// region Check if there is AuthToken ID in context
	identityIDVal := ctx.Value(AuthIdentityIDContextKey)
	if identityIDVal == nil {
		return nil, errors.NewUnauthorizedError(fmt.Errorf("wrong token 1"))
	}
	identityID, ok := identityIDVal.(int)
	if !ok {
		return nil, errors.NewUnauthorizedError(fmt.Errorf("wrong token 2"))
	}
	// endregion

	// region Check if there is JWTType in context
	jwtTypeVal := ctx.Value(JWTTypeContextKey)
	if jwtTypeVal == nil {
		return nil, errors.NewUnauthorizedError(fmt.Errorf("wrong token 3"))
	}
	jwtType, ok := jwtTypeVal.(model.JWTType)
	if !ok || jwtType != model.JWTTypeRefresh {
		return nil, errors.NewUnauthorizedError(fmt.Errorf("wrong token 4"))
	}
	// endregion

	i := model.AuthIdentity{ID: identityID}
	authJWT, err := r.store.JWTStore.MakeAuthJWT(&i)
	if err != nil {
		return nil, fmt.Errorf("can't make auth jwt for i %d: %w", i.ID, err)
	}
	authJWTStr, err := r.store.JWTStore.GetJWTString(authJWT)
	if err != nil {
		return nil, fmt.Errorf("can't make auth jwt string for i %d: %w", i.ID, err)
	}

	return &gql.GetAuthJWTResponse{
		AuthToken: *authJWTStr,
	}, nil
}
