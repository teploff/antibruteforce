package limiter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/domain/repository"
	bucket "github.com/teploff/antibruteforce/internal/implementation/repository"
	"github.com/teploff/antibruteforce/internal/shared"
)

const sleepTime = time.Millisecond * 100

type RateLimiter struct {
	loginBuckets    repository.BucketStorable
	passwordBuckets repository.BucketStorable
	ipBuckets       repository.BucketStorable
	ctx             context.Context
	duration        time.Duration
}

func NewRateLimiter(ctx context.Context, cfg config.RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		loginBuckets:    bucket.NewLeakyBucket(cfg.Login.RPM, cfg.Login.ExpireTime),
		passwordBuckets: bucket.NewLeakyBucket(cfg.Password.RPM, cfg.Password.ExpireTime),
		ipBuckets:       bucket.NewLeakyBucket(cfg.IP.RPM, cfg.IP.ExpireTime),
		ctx:             ctx,
		duration:        time.Duration(cfg.GCTime) * time.Second,
	}
}

func (r RateLimiter) IsBruteForce(login, password, ip string) (bool, error) {
	isAllowed, err := isRequestAllowed(r.loginBuckets, login)
	if err != nil || !isAllowed {
		return !isAllowed, err
	}

	isAllowed, err = isRequestAllowed(r.passwordBuckets, password)
	if err != nil || !isAllowed {
		return !isAllowed, err
	}

	isAllowed, err = isRequestAllowed(r.ipBuckets, ip)
	if err != nil || !isAllowed {
		return !isAllowed, err
	}

	return false, nil
}

func isRequestAllowed(rate repository.BucketStorable, keyBucket string) (bool, error) {
	limiter, err := rate.Get(keyBucket)
	if errors.Is(err, shared.ErrNotFound) {
		if limiter, err = rate.Add(keyBucket); err != nil {
			return false, err
		}
	}

	if !limiter.Allow() {
		return false, nil
	}

	return true, nil
}

// Every r.duration seconds check if the buckets expire?
func (r RateLimiter) RunGarbageCollector() {
	ticker := time.NewTicker(r.duration)

	for {
		select {
		case <-r.ctx.Done():
			fmt.Println("Close the chenal")
			ticker.Stop()

			return
		case <-ticker.C:
			r.loginBuckets.Clean()
			r.passwordBuckets.Clean()
			r.ipBuckets.Clean()
		default:
			time.Sleep(sleepTime)
		}
	}
}
