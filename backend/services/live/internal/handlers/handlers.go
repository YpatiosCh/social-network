package handlers

import (
	"context"
	"net/http"
	"social-network/shared/go/eventsub"
	middleware "social-network/shared/go/http-middleware"
	"social-network/shared/go/ratelimit"
	"time"
)

type Handlers struct {
	CacheService CacheService
	EventBox     eventsub.EventBox[int64, []byte]
}

type CacheService interface {
	IncrEx(ctx context.Context, key string, expSeconds int64) (int, error)
	SetStr(ctx context.Context, key string, value string, exp time.Duration) error
	GetStr(ctx context.Context, key string) (any, error)
	SetObj(ctx context.Context, key string, value any, exp time.Duration) error
	GetObj(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
	TestRedisConnection() error
}

func NewHandlers(serviceName string, CacheService CacheService, eventBox eventsub.EventBox[int64, []byte]) *http.ServeMux {
	handlers := Handlers{
		CacheService: CacheService,
		EventBox:     eventBox,
	}
	return handlers.BuildMux(serviceName)
}

// BuildMux builds and returns the HTTP request multiplexer with all routes and middleware applied
func (h *Handlers) BuildMux(serviceName string) *http.ServeMux {
	mux := http.NewServeMux()
	ratelimiter := ratelimit.NewRateLimiter(serviceName+":", h.CacheService)
	middlewareObj := middleware.NewMiddleware(ratelimiter, serviceName)
	Chain := middlewareObj.Chain

	IP := middleware.IPLimit
	USERID := middleware.UserLimit

	mux.HandleFunc("/test",
		Chain("/test").
			RateLimit(IP, 10, 10).
			AllowedMethod("GET").
			RateLimit(IP, 20, 5).
			EnrichContext().
			RateLimit(USERID, 10, 10).
			Finalize(h.testHandler()))

	mux.HandleFunc("/live",
		Chain("/live").
			RateLimit(IP, 10, 10).
			AllowedMethod("GET").
			Auth().
			EnrichContext().
			RateLimit(USERID, 10, 10).
			Finalize(h.Connect()))

	return mux
}
