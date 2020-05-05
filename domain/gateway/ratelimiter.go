package gateway

import "github.com/teploff/antibruteforce/domain/entity"

type RateLimiter interface {
	AddBucket(key string) (*entity.Limiter, error)
	DeleteBucket(key string) error
	GetLimiter(key string) (*entity.Limiter, error)
}
