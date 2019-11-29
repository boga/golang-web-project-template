package resolvers

import (
	"cipherassets.core/gql"
	"cipherassets.core/store"
)

type Resolver struct {
	store *store.Store
}

func NewResolver(store *store.Store) *Resolver {
	return &Resolver{store: store}
}

func (r *Resolver) Mutation() gql.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() gql.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
