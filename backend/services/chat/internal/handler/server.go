package handler

import (
	"social-network/services/chat/internal/application"

	_ "github.com/lib/pq"
)

type ChatHandler struct {
	// pb.UnimplementedChatServiceServer
	Application *application.ChatService
	Port        string
}
