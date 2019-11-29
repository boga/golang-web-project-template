package resolvers

import (
	"context"

	"cipherassets.core/model"
)

type userResolver struct{ *Resolver }

func (r *userResolver) Banned(ctx context.Context, obj *model.User) (bool, error) {
	return false, nil
}

// func (r *userResolver) AuthIdentities(ctx context.Context, user *model.User) ([]*model.AuthIdentity, error) {
// 	return []*model.AuthIdentity{
// 		{
// 			ID:  1,
// 			UID: "First identity",
// 		},
// 		{
// 			ID:  2,
// 			UID: "Second identity",
// 		},
// 	}, nil
// }
