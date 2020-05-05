package gateway

import (
	"sync"

	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/gateway"
	"github.com/teploff/antibruteforce/internal/shared"
)

type rateLimiter struct {
	buckets map[string]*entity.Limiter
	rate    entity.Limit
	mu      *sync.RWMutex
}

func NewRateLimiter(rate int) gateway.RateLimiter {
	return &rateLimiter{
		buckets: make(map[string]*entity.Limiter),
		rate:    entity.Limit(rate),
		mu:      &sync.RWMutex{},
	}
}

func (r *rateLimiter) AddBucket(key string) (*entity.Limiter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.buckets[key]; exist {
		return nil, shared.ErrAlreadyExist
	}

	r.buckets[key] = entity.NewLimiter(r.rate, 3)

	return r.buckets[key], nil
}

func (r *rateLimiter) DeleteBucket(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.buckets[key]; !exist {
		return shared.ErrNotFound
	}

	delete(r.buckets, key)

	return nil
}

func (r *rateLimiter) GetLimiter(key string) (*entity.Limiter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	limiter, exist := r.buckets[key]
	if !exist {
		return nil, shared.ErrNotFound
	}

	return limiter, nil
}
