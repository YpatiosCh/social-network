/*
Expose methods via gRpc
*/

package server

import (
	"context"
	"fmt"
	"social-network/services/users/internal/application"
	pb "social-network/shared/gen-go/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) GetBasicUserInfo(ctx context.Context, req *wrapperspb.Int64Value) (*pb.User, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userId := req.GetValue()
	if userId == 0 {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	u, err := s.Service.GetBasicUserInfo(ctx, req.GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetBasicUserInfo: %v", err)
	}

	return &pb.User{
		UserId:   u.UserId,
		Username: u.Username,
		Avatar:   u.Avatar,
	}, nil
}

func (s *Server) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfileResponse, error) {
	fmt.Println("GetUserProfile gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userID := req.GetUserId()
	reqID := req.GetRequesterId()
	if userID == 0 || reqID == 0 {
		return nil, status.Error(codes.InvalidArgument, "GetUserProfile: userId and requesterId are required")
	}
	userProfileRequest := application.UserProfileRequest{
		UserId:      req.GetUserId(),
		RequesterId: req.GetRequesterId(),
	}

	profile, err := s.Service.GetUserProfile(ctx, userProfileRequest)
	if err != nil {
		fmt.Println("Error in GetUserProfile:", err)
		return nil, status.Errorf(codes.Internal, "GetUserProfile: %v", err)
	}

	return &pb.UserProfileResponse{
		UserId:   profile.UserId,
		Username: profile.Username,
	}, nil
}

func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.User, error) {
	fmt.Println("RegisterUser gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	user, err := s.Service.RegisterUser(ctx, application.RegisterUserRequest{
		Username:    req.GetUsername(),
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		DateOfBirth: req.GetDateOfBirth().AsTime(),
		Avatar:      req.GetAvatar(),
		About:       req.GetAbout(),
		Public:      req.GetPublic(),
	})
	if err != nil {
		fmt.Println("Error in RegisterUser:", err)
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.User{
		UserId:   user.UserId,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.User, error) {
	fmt.Println("LoginUser gRPC method called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "LoginUser: request is nil")
	}

	ident := req.GetIdentifier()
	pass := req.GetPassword()
	if ident == "" || pass == "" {
		return nil, status.Error(codes.InvalidArgument, "LoginUser: identifier and password are required")
	}

	user, err := s.Service.LoginUser(ctx, application.LoginRequest{
		Identifier: ident,
		Password:   pass,
	})
	if err != nil {
		fmt.Println("Error in LoginUser:", err)
		return nil, status.Errorf(codes.Internal, "LoginUser: failed to login user: %v", err)
	}

	return &pb.User{
		UserId:   user.UserId,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

func (s *Server) UpdateUserPassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "UpdateUserPassword: request is nil")
	}

	userID := req.GetUserId()
	newPassword := req.GetNewPassword()
	if userID == 0 || newPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "UpdateUserPassword: userId and newPassword are required")
	}

	err := s.Service.UpdateUserPassword(ctx, application.UpdatePasswordRequest{
		UserId:      userID,
		NewPassword: newPassword,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UpdateUserPassword: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateUserEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "UpdateUserEmail: request is nil")
	}

	userID := req.GetUserId()
	newEmail := req.GetEmail()
	if userID == 0 || newEmail == "" {
		return nil, status.Error(codes.InvalidArgument, "UpdateUserEmail: userId and newEmail are required")
	}

	err := s.Service.UpdateUserEmail(ctx, application.UpdateEmailRequest{
		UserId: userID,
		Email:  newEmail,
	})
	if err != nil {
		fmt.Println("Error in UpdateUserEmail:", err)
		return nil, status.Errorf(codes.Internal, "UpdateUserEmail: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetFollowersPaginated(ctx context.Context, req *pb.Pagination) (*pb.ListUsers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetFollowersPaginated: request is nil")
	}

	userId := req.GetUserId()
	limit := req.GetLimit()
	offset := req.GetOffset()
	if userId == 0 || offset < 0 || limit > application.MAX_FOLLOWERS_PAGE_LIMIT {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("UpdateUserEmail: userId %v, limit: %v, offset: %v", userId, limit, offset))
	}

	pag := application.Pagination{
		UserId: userId,
		Limit:  limit,
		Offset: offset,
	}

	resp, err := s.Service.GetFollowersPaginated(ctx, pag)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetFollowersPaginated: %v", err)
	}
	return usersToPB(resp), nil
}

func (s *Server) GetFollowingPaginated(ctx context.Context, req *pb.Pagination) (*pb.ListUsers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetFollowingPaginated: request is nil")
	}

	userId := req.GetUserId()
	limit := req.GetLimit()
	offset := req.GetOffset()
	if userId == 0 || offset < 0 || limit > application.MAX_FOLLOWERS_PAGE_LIMIT {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("GetFollowingPaginated: userId %v, limit: %v, offset: %v", userId, limit, offset))
	}

	pag := application.Pagination{
		UserId: userId,
		Limit:  limit,
		Offset: offset,
	}

	resp, err := s.Service.GetFollowingPaginated(ctx, pag)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetFollowingPaginated: %v", err)
	}
	return usersToPB(resp), nil
}

func (s *Server) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "FollowUser: request is nil")
	}

	followerId := req.GetFollowerId()
	targetUserId := req.GetTargetUserId()
	if followerId == 0 || targetUserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "follower id and target user id are required")
	}

	resp, err := s.Service.FollowUser(ctx, application.FollowUserReq{
		FollowerId:   followerId,
		TargetUserId: targetUserId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "FollowUser: %v", err)
	}

	return &pb.FollowUserResponse{
		IsPending:         resp.IsPending,
		ViewerIsFollowing: resp.ViewerIsFollowing,
	}, nil
}

func (s *Server) UnFollowUser(ctx context.Context, req *pb.FollowUserRequest) (*wrapperspb.BoolValue, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "UnFollowUser: request is nil")
	}
	followerId := req.GetFollowerId()
	targetId := req.GetTargetUserId()
	if followerId == 0 || targetId == 0 {
		return nil, status.Error(codes.InvalidArgument, "follower id and target user id are required")
	}

	resp, err := s.Service.UnFollowUser(ctx, application.FollowUserReq{
		FollowerId:   followerId,
		TargetUserId: targetId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UnFollowUser: %v", err)
	}

	return wrapperspb.Bool(resp), nil
}

func (s *Server) HandleFollowRequest(ctx context.Context, req *pb.HandleFollowRequestRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "HandleFollowRequest: request is nil")
	}

	userID := req.GetUserId()
	reqID := req.GetRequesterId()
	acc := req.GetAccept()
	if userID == 0 || reqID == 0 {
		return nil, status.Error(codes.InvalidArgument, "follower id and requester id are required")
	}
	err := s.Service.HandleFollowRequest(ctx, application.HandleFollowRequestReq{
		UserId:      userID,
		RequesterId: reqID,
		Accept:      acc,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "HandleFollowRequest: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetFollowingIds(ctx context.Context, req *wrapperspb.Int64Value) (*pb.Int64Arr, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetFollowingIds: request is nil")
	}

	resp, err := s.Service.GetFollowingIds(ctx, req.GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetFollowingIds: %v", err)
	}
	return &pb.Int64Arr{Values: resp}, nil
}

func (s *Server) GetFollowSuggestions(ctx context.Context, req *wrapperspb.Int64Value) (*pb.ListUsers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "GetFollowSuggestions: request is nil")
	}

	resp, err := s.Service.GetFollowSuggestions(ctx, req.GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetFollowSuggestions: %v", err)
	}

	return usersToPB(resp), nil
}

func usersToPB(dbUsers []application.User) *pb.ListUsers {
	pbUsers := make([]*pb.User, 0, len(dbUsers))

	for _, u := range dbUsers {
		pbUsers = append(pbUsers, &pb.User{
			UserId:   u.UserId,
			Username: u.Username,
			Avatar:   u.Avatar,
		})
	}

	return &pb.ListUsers{Users: pbUsers}
}
