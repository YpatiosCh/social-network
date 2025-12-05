package handler

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"social-network/services/chat/internal/application"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type ChatHandler struct {
	// pb.UnimplementedChatServiceServer
	Application *application.ChatService
	Port        string
}

func Run() error {
	app, err := application.Run(context.Background())
	if err != nil {
		return err
	}

	service := &ChatHandler{
		Application: app,
		Port:        ":50051",
	}

	log.Println("Running gRpc service...")
	grpc := RunGRPCServer(service)

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	grpc.GracefulStop()
	log.Println("Server stopped")
	return nil
}

// RunGRPCServer starts the gRPC server and blocks
func RunGRPCServer(s *ChatHandler) *grpc.Server {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Port, err)
	}

	grpcServer := grpc.NewServer()

	// pb.RegisterChatServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", s.Port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	return grpcServer
}
