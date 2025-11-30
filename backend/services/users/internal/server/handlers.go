/*
Expose methods via gRpc
*/

package server

import (
	"context"
	"fmt"
	"social-network/services/users/internal/application"
	pb "social-network/shared/gen-go/users"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) GetBasicUserInfo(ctx context.Context, req *wrapperspb.Int64Value) (*pb.User, error) {
	u, err := s.Service.GetBasicUserInfo(ctx, req.GetValue())
	return &pb.User{
		UserId:   u.UserId,
		Username: u.Username,
		Avatar:   u.Avatar,
	}, err
}

func (s *Server) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfileResponse, error) {
	fmt.Println("GetUserProfile gRPC method called")
	userProfileRequest := application.UserProfileRequest{
		UserId:      req.GetUserId(),
		RequesterId: req.GetRequesterId(),
	}

	profile, err := s.Service.GetUserProfile(ctx, userProfileRequest)
	if err != nil {
		fmt.Println("Error in GetUserProfile:", err)
		return nil, err
	}

	return &pb.UserProfileResponse{
		UserId:   profile.UserId,
		Username: profile.Username,
	}, nil
}
