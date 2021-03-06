package bucket

import (
	"sync"
	"time"

	"github.com/teploff/antibruteforce/internal/domain/entity"
	"github.com/teploff/antibruteforce/internal/domain/repository"
	"github.com/teploff/antibruteforce/internal/shared"
)

type leakyBucket struct {
	buckets    map[string]*entity.Limiter
	rate       int
	interval   time.Duration
	mu         *sync.RWMutex
	expireTime time.Duration
}

// NewLeakyBucket returns instance of Leaky buckets.
func NewLeakyBucket(rate int, interval, expireTime time.Duration) repository.BucketStorable {
	return &leakyBucket{
		buckets:    make(map[string]*entity.Limiter),
		rate:       rate,
		interval:   interval,
		mu:         &sync.RWMutex{},
		expireTime: expireTime,
	}
}

func (l *leakyBucket) Add(key string) (*entity.Limiter, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exist := l.buckets[key]; exist {
		return nil, shared.ErrAlreadyExist
	}

	l.buckets[key] = entity.NewLimiter(l.rate, l.interval)

	return l.buckets[key], nil
}

func (l *leakyBucket) Delete(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exist := l.buckets[key]; !exist {
		return shared.ErrNotFound
	}

	delete(l.buckets, key)

	return nil
}

func (l *leakyBucket) Get(key string) (*entity.Limiter, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	limiter, exist := l.buckets[key]
	if !exist {
		return nil, shared.ErrNotFound
	}

	return limiter, nil
}

func (l *leakyBucket) Clean() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for k, v := range l.buckets {
		if time.Since(v.LastSeen()) > l.expireTime {
			delete(l.buckets, k)
		}
	}
}
