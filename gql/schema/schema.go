package schema

import (
	"github.com/99designs/gqlgen/graphql"

	core "cipherassets.core"
	"cipherassets.core/gql"
	"cipherassets.core/gql/directives"
	"cipherassets.core/gql/resolvers"
	"cipherassets.core/store"
)

func NewSchema(config *core.Config, store *store.Store) graphql.ExecutableSchema {
	c := gql.Config{
		Directives: directives.NewDirectives(store),
		Resolvers:  resolvers.NewResolver(config, store),
	}
	// c.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver, addUserToCtx *bool) (interface{}, error) {
	// 	log.Printf("auth 1 directive %+v", addUserToCtx)
	// 	_ = store
	// 	return next(ctx)
	// }
	// c.Complexity.Mutation.TotpGenerate = func(childComplexity int) int {
	// 	log.Printf("auth 1 directive")
	// 	return 0
	// }

	return gql.NewExecutableSchema(c)
}
