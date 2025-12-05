package remoteservices

import (
	"fmt"
	"log"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRpcServices struct {
	Users users.UserServiceClient
	Chat  chat.ChatServiceClient
}

func NewServices() GRpcServices {
	return GRpcServices{}
}

func (g *GRpcServices) StartConnections() (func(), error) {
	usersConn, err := grpc.NewClient("users:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(usersConn.CanonicalTarget())
	fmt.Println(usersConn.GetState())
	fmt.Println(usersConn.Target())
	g.Users = users.NewUserServiceClient(usersConn)

	chatConn, err := grpc.NewClient("chat:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	g.Chat = chat.NewChatServiceClient(chatConn)

	deferMe := func() {
		usersConn.Close()
		chatConn.Close()
	}
	return deferMe, nil
}
