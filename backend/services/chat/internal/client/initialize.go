package client

import (
	userpb "social-network/shared/gen-go/users"
)

type Clients struct {
	UserClient userpb.UserServiceClient
}
