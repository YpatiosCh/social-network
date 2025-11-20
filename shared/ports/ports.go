package ports

// Container ports should use kubernetes DNS when up
const (
	Users         string = "users:50051"
	Forum         string = "forum:50052"
	Chat          string = "chat:"
	Notifications string = "notifications:"
)
