package repository

import (
	"sync"
	"time"

	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/repository"
	"github.com/teploff/antibruteforce/internal/shared"
)

type leakyBucket struct {
	buckets    map[string]*entity.Limiter
	rate       entity.Limit
	mu         *sync.RWMutex
	expireTime time.Duration
}

func NewLeakyBucket(rate, expireTime int) repository.BucketStorable {
	return &leakyBucket{
		buckets:    make(map[string]*entity.Limiter),
		rate:       entity.Limit(rate),
		mu:         &sync.RWMutex{},
		expireTime: time.Duration(expireTime) * time.Second,
	}
}

func (l *leakyBucket) Add(key string) (*entity.Limiter, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exist := l.buckets[key]; exist {
		return nil, shared.ErrAlreadyExist
	}

	l.buckets[key] = entity.NewLimiter(l.rate, 3)

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
	l.mu.Lock()
	defer l.mu.Unlock()

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
		if time.Since(v.Last()) > l.expireTime {
			delete(l.buckets, k)
		}
	}
}
