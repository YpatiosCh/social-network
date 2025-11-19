package server

import (
	"fmt"
	userpb "social-network/shared/gen/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Server) InitClients() {
	if err := s.InitUserClient(); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) InitUserClient() error {
	conn, err := grpc.NewClient("users:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to dial user service: %v", err)
	}
	s.clients.UserClient = userpb.NewUserServiceClient(conn)
	return nil
}
