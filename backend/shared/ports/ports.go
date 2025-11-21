/*
This is a namespace to hold port addreses for all services.
*/
package ports

// Container ports used with individual docker containers
const (
	Users         string = "users:50051"
	Forum         string = "forum:50052"
	Chat          string = "chat:"
	Notifications string = "notifications:"
	Example       string = "example:1234"
)

// Kube service DNS
const (
	UsersDNS string = "users.social-network.svc.cluster.local:50051"
	UsersDb  string = "users-db.social-network.svc.cluster.local:5432"
)
