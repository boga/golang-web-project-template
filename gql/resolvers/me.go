package resolvers

import (
	"context"

	"cipherassets.core/model"
)

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	name := "Chuck Norris"
	us, err := r.store.UserStore.GetUsers()
	_, _ = us, err
	return &model.User{
		ID:             1,
		Name:           &name,
		AuthIdentities: nil,
	}, nil
}
