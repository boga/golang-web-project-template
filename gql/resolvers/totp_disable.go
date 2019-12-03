package resolvers

import (
	"context"
	"fmt"

	"cipherassets.core/gql"
	"cipherassets.core/gql/errors"
	"cipherassets.core/model"
)

func (r *mutationResolver) TotpDisable(ctx context.Context, data gql.TOTPDisableInput) (*gql.TOTPDisableResponse, error) {
	i := ctx.Value(AuthIdentityContextKey).(*model.AuthIdentity)
	user, err := r.store.UserStore.FindUserByID(i.UserID)
	if err != nil {
		return nil, err
	}

	if !user.TOTPEnabled {
		return nil, errors.NewTOTPNotEnabledError(nil)
	}

	err = r.store.UserStore.ValidateTOTPCode(user, data.Code, false)
	if err != nil {
		return nil, errors.NewTOTPNotValidError(err)
	}

	user.TOTPEnabled = false
	user.TOTPSecret = nil
	if err := r.store.UserStore.SaveUser(user); err != nil {
		return nil, fmt.Errorf("can't save TOTP disabled to user (%d): %w", user.ID, err)
	}

	return &gql.TOTPDisableResponse{
		Success: true,
	}, nil
}
