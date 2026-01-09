package ct

import "fmt"

func PrivateMessageKey(receiverId any) string {
	return fmt.Sprintf("dm.%v", receiverId)
}

func GroupMessageKey(conversationId any) string {
	return fmt.Sprintf("grm.%v", conversationId)
}

func NotificationKey(receiverId any) string {
	return fmt.Sprintf("ntf.%v", receiverId)
}
