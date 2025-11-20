package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	commonpb "social-network/shared/gen/common"
	pb "social-network/shared/gen/users"
	"social-network/shared/ports"

	"google.golang.org/grpc"
)

// Server struct placeholder
type Server struct {
	pb.UnimplementedUserServiceServer
	clients Clients
	Port    string
}

type Clients struct {
	UserClient pb.UserServiceClient
}

// RunGRPCServer starts the gRPC server and blocks
func (s *Server) RunGRPCServer() {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	// pb.RegisterUserServiceServer(grpcServer, &Server{})

	// ================================================
	// TEST GRPC CONN
	// ================================================

	// Wait for other servers to run first
	time.Sleep(time.Second)
	s.InitClients()

	// Example call to UserClient
	resp, err := s.clients.UserClient.GetBasicUserInfo(context.Background(), &commonpb.UserId{
		Id: 1234,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.String())

	// -------------------------------------------------

	log.Printf("gRPC server listening on %s", s.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func NewForumServer() *Server {
	return &Server{
		Port: ports.Forum,
	}
}
