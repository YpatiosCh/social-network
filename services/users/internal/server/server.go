package server

import (
	"log"
	"net"

	userservice "social-network/services/users/internal/service"
	pb "social-network/shared/gen/users"

	"google.golang.org/grpc"
)

// Server struct placeholder
type Server struct {
	// Add DB or other dependencies here later
	pb.UnimplementedUserServiceServer
	clients Clients
	Port    string
	Service *userservice.UserService
}

type Clients struct {
	Example pb.UserServiceClient
}

// RunGRPCServer starts the gRPC server and blocks
func (s *Server) RunGRPCServer() {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	// TODO: Register services here, e.g.,
	pb.RegisterUserServiceServer(grpcServer, &Server{})

	log.Printf("gRPC server listening on %s", s.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func NewUsersServer(port string, service *userservice.UserService) *Server {
	return &Server{
		Port:    port,
		clients: Clients{},
		Service: service,
	}
}
