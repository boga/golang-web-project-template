package resolvers

import (
	core "cipherassets.core"
	"cipherassets.core/gql"
	"cipherassets.core/store"
)

type ContextKey string

var AuthIdentityIDContextKey = ContextKey("AuthIdentityID")
var AuthIdentityContextKey = ContextKey("AuthIdentity")

type Resolver struct {
	config *core.Config
	store  *store.Store
}

func NewResolver(config *core.Config, store *store.Store) *Resolver {
	return &Resolver{
		config: config,
		store:  store,
	}
}

func (r *Resolver) Mutation() gql.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() gql.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
