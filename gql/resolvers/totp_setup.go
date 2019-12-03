package resolvers

import (
	"context"
	"strings"

	"cipherassets.core/gql"
	"cipherassets.core/model"
)

func (r *mutationResolver) TotpSetup(ctx context.Context, data gql.TOTPSetupInput) (*gql.TOTPSetupResponse, error) {
	identity := ctx.Value(AuthIdentityContextKey).(*model.AuthIdentity)
	user, err := r.store.UserStore.FindUserByID(identity.UserID)
	if err != nil {
		return nil, err
	}

	err = r.store.UserStore.ValidateTOTPCode(user, data.Code, true)
	if err != nil {
		return nil, err
	}

	user.TOTPEnabled = true
	backupCodesArr := r.store.UserStore.NewBackupCodes()
	backupCodes := strings.Join(backupCodesArr, "\n")
	user.TOTPBackupCodes = &backupCodes
	if err := r.store.UserStore.SaveUser(user); err != nil {
		return nil, err
	}

	return &gql.TOTPSetupResponse{
		BackupCodes: backupCodesArr,
	}, nil
}
