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
