package resolvers

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"net/http"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"cipherassets.core/gql"
	"cipherassets.core/gql/errors"
	"cipherassets.core/model"
)

func (r *mutationResolver) TotpGenerate(ctx context.Context) (*gql.TOTPGenerateResponse, error) {
	identity := ctx.Value(AuthIdentityContextKey).(*model.AuthIdentity)
	user, err := r.store.UserStore.FindUserByID(identity.UserID)
	if err != nil {
		return nil, err
	}

	if user.TOTPEnabled {
		return nil, errors.NewApiError(nil, "totp-enabled", "totp is already enabled for user", http.StatusBadRequest)
	}

	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      r.config.App.Name,
		AccountName: identity.UID,
		// Period:      60, // seconds
		// SecretSize:  0,
		// Secret:      nil,
		// Digits:    6,
		Algorithm: otp.AlgorithmSHA1,
	})
	secret := key.Secret()
	image, err := key.Image(600, 600)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	if err := png.Encode(&buff, image); err != nil {
		return nil, err
	}
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())

	user.TOTPSecret = &secret
	if err := r.store.UserStore.SaveUser(user); err != nil {
		return nil, err
	}

	return &gql.TOTPGenerateResponse{
		Qrcode: encodedString,
	}, nil
}
