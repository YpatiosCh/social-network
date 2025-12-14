package retrieveusers

import (
	"context"
	"fmt"

	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"
	redis_connector "social-network/shared/go/redis"
	"time"
)

type UserRetriever struct {
	clients UsersBatchClient
	cache   RedisCache
	ttl     time.Duration
}

func NewUserRetriever(clients UsersBatchClient, cache *redis_connector.RedisClient, ttl time.Duration) *UserRetriever {
	return &UserRetriever{clients: clients, cache: cache, ttl: ttl}
}

// GetUsers returns a map[userID]User, using cache + batch RPC.
func (h *UserRetriever) GetUsers(ctx context.Context, userIDs []int64) (map[int64]models.User, error) {
	idSet := make(map[int64]struct{}, len(userIDs))
	for _, id := range userIDs {
		idSet[id] = struct{}{}
	}

	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	users := make(map[int64]models.User, len(ids))
	var missing []int64

	// Redis lookup
	for _, id := range ids {
		var u models.User
		if err := h.cache.GetObj(ctx, fmt.Sprintf("basic_user_info:%d", id), &u); err == nil {
			users[id] = u
		} else {
			missing = append(missing, id)
		}
	}

	// Batch RPC for missing users
	if len(missing) > 0 {
		resp, err := h.clients.GetBatchBasicUserInfo(ctx, missing)
		if err != nil {
			return nil, err
		}

		for _, u := range resp.Users {
			user := models.User{
				UserId:   ct.Id(u.UserId),
				Username: ct.Username(u.Username),
				AvatarId: ct.Id(u.Avatar),
			}
			users[u.UserId] = user
			_ = h.cache.SetObj(ctx,
				fmt.Sprintf("basic_user_info:%d", u.UserId),
				user,
				h.ttl,
			)
		}
	}

	return users, nil
}
