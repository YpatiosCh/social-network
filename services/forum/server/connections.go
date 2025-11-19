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
	// List of initializer functions
	initializers := []func() error{
		s.InitUserClient,
		// Add more here as you add more clients
	}

	for _, initFn := range initializers {
		if err := initFn(); err != nil {
			fmt.Println(err)
		}
	}
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
