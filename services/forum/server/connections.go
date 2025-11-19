/*
Establish connections to other services
*/

package server

import (
	"fmt"
	userpb "social-network/shared/gen/users"
	"social-network/shared/ports"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Initialize connections to clients
func (s *Server) InitClients() {
	if err := s.InitUserClient(); err != nil {
		fmt.Println(err)
	}
	// Add more clients
}

// Connects to client and adds connection to s.Clients
func (s *Server) InitUserClient() error {
	conn, err := grpc.NewClient(ports.Users, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to dial user service: %v", err)
	}
	s.clients.UserClient = userpb.NewUserServiceClient(conn)
	return nil
}
