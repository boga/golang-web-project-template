package resolvers

import (
	"context"

	"cipherassets.core/gql"
	"cipherassets.core/model"
)

func (r *mutationResolver) CreateAuthIdentity(ctx context.Context, input gql.NewAuthIdentity) (*model.AuthIdentity, error) {
	name := "Chuck Norris"
	return &model.AuthIdentity{
		ID:  1,
		UID: "some email",
		User: model.User{
			ID:             1,
			Name:           &name,
			AuthIdentities: nil,
		},
	}, nil
}
