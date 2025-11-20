/*
Expose methods via gRpc
*/

package server

import (
	"context"
	userservice "social-network/services/users/internal/service"
	commonpb "social-network/shared/gen/common"
	pb "social-network/shared/gen/users"
)

func (s *Server) GetBasicUserInfo(ctx context.Context, req *commonpb.UserId) (*pb.BasicUserInfo, error) {
	u, err := userservice.GetBasicUserInfo(ctx, req.Id)
	return &pb.BasicUserInfo{
		UserName:      u.UserName,
		Avatar:        u.Avatar,
		PublicProfile: u.PublicProfile,
	}, err
}
