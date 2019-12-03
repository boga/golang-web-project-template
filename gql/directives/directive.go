package directives

import (
	"cipherassets.core/gql"
	"cipherassets.core/store"
)

func NewDirectives(store *store.Store) gql.DirectiveRoot {
	return gql.DirectiveRoot{
		Auth: NewAuthDirective(store),
	}
}
