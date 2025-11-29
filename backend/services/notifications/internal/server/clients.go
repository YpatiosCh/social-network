package server

import (
	"fmt"
	"social-network/shared/ports"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

// Initialize connections to clients. Each one is called with dial options
func (s *Server) InitClients() {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
        "loadBalancingConfig": [{"round_robin":{}}]
    	}`),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 2 * time.Second,
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.2,
				Jitter:     0.5,
				MaxDelay:   5 * time.Second,
			},
		}),
	}

	// List of initializer functions
	initializers := []func(opts []grpc.DialOption) error{
		// Add all init clients funcs
		s.InitTemplateClient,
	}

	for _, initFn := range initializers {
		if err := initFn(dialOpts); err != nil {
			fmt.Println(err)
		}
	}
}

// Connects to client and adds connection to s.Clients
// TEMPLATE
func (s *Server) InitTemplateClient(opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(ports.Users, opts...)
	if err != nil {
		err = fmt.Errorf("failed to dial user service: %v", err)
	}
	_ = conn
	// s.Clients.Example = explPb.NewUserServiceClient(conn)
	return err
}
