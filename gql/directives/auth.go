package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"

	"cipherassets.core/gql/errors"
	"cipherassets.core/gql/resolvers"
	"cipherassets.core/store"
)

func NewAuthDirective(s *store.Store) func(ctx context.Context, obj interface{}, next graphql.Resolver, addUserToCtx *bool) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, addUserToCtx *bool) (res interface{}, err error) {
		identityIDVal := ctx.Value(resolvers.AuthIdentityIDContextKey)
		identityID, ok := identityIDVal.(int)
		if !ok {
			return nil, errors.NewUnauthorizedError(fmt.Errorf("wrong token"))
		}

		if (addUserToCtx == nil) || !*addUserToCtx {
			return next(ctx)
		}

		identity, err := s.UserStore.FindIdentityByID(identityID)
		if err != nil {
			return nil, errors.NewUnauthorizedError(err)
		}

		ctx = context.WithValue(ctx, resolvers.AuthIdentityContextKey, identity)

		return next(ctx)
	}
}
