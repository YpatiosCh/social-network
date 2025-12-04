package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	// pb.UnimplementedChatServiceServer
	Clients Clients
	Port    string
	// Application *us.Application
}

// Holds connections to clients
type Clients struct {
}

// RunGRPCServer starts the gRPC server and blocks
func (s *Server) RunGRPCServer() {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	// pb.RegisterChatServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", s.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
