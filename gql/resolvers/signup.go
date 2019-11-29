package resolvers

import (
	"context"
	"fmt"

	"cipherassets.core/gql"
)

func (r *mutationResolver) Signup(ctx context.Context, creds gql.SignupInput) (*gql.SignupResponse, error) {
	i, err := r.store.UserStore.FindIdentityByUID(creds.Email)
	if err != nil {
		return nil, fmt.Errorf("identity not found: %w", err)
	}
	if i != nil {
		return nil, fmt.Errorf("email '%s' taken", creds.Email)
	}
	// region Creating user
	u, err := r.store.UserStore.NewUser()
	if err != nil {
		return nil, fmt.Errorf("can't create user for %s: %w", creds.Email, err)
	}
	err = r.store.UserStore.SaveUser(u)
	if err != nil {
		return nil, fmt.Errorf("can't save user for %s: %w", creds.Email, err)
	}
	// endregion

	// region Creating identity
	i, err = r.store.UserStore.NewIdentity(creds.Email, creds.Password, u.ID)
	if err != nil {
		return nil, fmt.Errorf("can't create identity for %s: %w", creds.Email, err)
	}
	err = r.store.UserStore.SaveIdentity(i)
	if err != nil {
		return nil, fmt.Errorf("can't save identity for %s: %w", creds.Email, err)
	}
	// endregion

	return &gql.SignupResponse{
		User: u,
	}, nil
}
