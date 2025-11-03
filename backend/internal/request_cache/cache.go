/*
Creates and manages a cache to store request IDs to prevent duplicate processing of requests.
Usage:
1. Initialize the cache with a specified capacity using New(capacity int).
2. Use RequestExists(requestId string) to check if a request ID already exists in the cache.
3. Use AddReq on success

OPTIONAL:
Get previous request by userId
*/
package requestcache

import (
	"sync"
)

// TODO low priority improve efficiency of removing requests from map, and
type RequestMeta struct {
	Timestamp string
	UserId    int64
}

type RequestCache struct {
	mu         sync.RWMutex
	requestIds map[string]RequestMeta
	order      []string
	capacity   int
}

func New(capacity int) *RequestCache {
	return &RequestCache{
		requestIds: make(map[string]RequestMeta),
		order:      make([]string, 0),
		capacity:   capacity,
	}
}

func (rc *RequestCache) RequestExists(requestId string) (RequestMeta, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	req, exists := rc.requestIds[requestId]
	return req, exists
}

func (rc *RequestCache) AddReq(req RequestMeta, requestId string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if _, exists := rc.requestIds[requestId]; exists {
		rc.requestIds[requestId] = req
		return
	}

	if len(rc.requestIds) >= rc.capacity {
		oldestID := rc.order[0]
		delete(rc.requestIds, oldestID)
		rc.order = rc.order[1:]
	}

	rc.requestIds[requestId] = req
	rc.order = append(rc.order, requestId)
}

func (rc *RequestCache) LastUserReq(userId int64) (RequestMeta, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	for i := len(rc.order) - 1; i >= 0; i-- {
		reqID := rc.order[i]
		meta := rc.requestIds[reqID]
		if meta.UserId == userId {
			return meta, true
		}
	}

	return RequestMeta{}, false
}
