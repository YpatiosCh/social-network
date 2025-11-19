/*
Establish connections to other services
*/

package server

import "fmt"

// Initialize connections to clients
func (s *Server) InitClients() {
	// List of initializer functions
	initializers := []func() error{
		s.InitClient,
	}

	for _, initFn := range initializers {
		if err := initFn(); err != nil {
			fmt.Println(err)
		}
	}
}

// Connects to client and adds connection to s.Clients
func (s *Server) InitClient() error {
	return nil
}
