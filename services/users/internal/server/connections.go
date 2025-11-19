/*
Establish connections to other services
*/

package server

import "fmt"

// Initialize connections to clients
func (s *Server) InitClients() {
	if err := s.InitClient(); err != nil {
		fmt.Println(err)
	}
	// Add more clients
}

// Connects to client and adds connection to s.Clients
func (s *Server) InitClient() error {
	return nil
}
