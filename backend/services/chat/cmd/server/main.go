package main

import (
	"fmt"
	"social-network/services/chat/internal/handler"
)

func main() {
	if err := handler.Run(); err != nil {
		fmt.Println(err)
	}
}
