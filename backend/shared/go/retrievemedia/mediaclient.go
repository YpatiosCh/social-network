package retrievemedia

import (
	"context"
	"social-network/shared/gen-go/media"
	"time"
)

type Client interface {
	GetImages(ctx context.Context, req *media.GetImagesRequest, variant *media.FileVariant) (*media.GetImagesResponse, error)
	GetImage(ctx context.Context, req *media.GetImageRequest) (*media.GetImageResponse, error)
}

// RedisCache defines the minimal Redis operations used by the hydrator.
type RedisCache interface {
	GetStr(ctx context.Context, key string) (any, error)
	SetStr(ctx context.Context, key string, value string, exp time.Duration) error
}
