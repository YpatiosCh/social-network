package mapping

import (
	commonpb "social-network/shared/gen-go/common"
	md "social-network/shared/go/models"
)

func MapUserToProto(u md.User) *commonpb.User {
	return &commonpb.User{
		UserId:    u.UserId.Int64(),
		Username:  u.Username.String(),
		Avatar:    u.AvatarId.Int64(),
		AvatarUrl: u.AvatarURL,
	}
}
