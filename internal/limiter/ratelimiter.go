package limiter

import (
	"errors"
	"time"

	"github.com/teploff/antibruteforce/internal/domain/repository"
	"github.com/teploff/antibruteforce/internal/shared"
)

const sleepTime = time.Millisecond * 100

// RateLimiter is facade of three type leaky Buckets instance: login, password and ip.
type RateLimiter struct {
	loginBuckets    repository.BucketStorable
	passwordBuckets repository.BucketStorable
	ipBuckets       repository.BucketStorable
	duration        time.Duration
	cancelCh        chan struct{}
}

// NewRateLimiter returns rate limiter of instance.
func NewRateLimiter(login, password, ip repository.BucketStorable, d time.Duration) *RateLimiter {
	return &RateLimiter{
		loginBuckets:    login,
		passwordBuckets: password,
		ipBuckets:       ip,
		duration:        d,
		cancelCh:        make(chan struct{}),
	}
}

// IsBruteForce checking of brute-force attack.
func (r *RateLimiter) IsBruteForce(login, password, ip string) (bool, error) {
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

// RunGarbageCollector run every r.duration time to delete expired buckets.
func (r *RateLimiter) RunGarbageCollector() {
	ticker := time.NewTicker(r.duration)

	for {
		select {
		case <-r.cancelCh:
			ticker.Stop()
			close(r.cancelCh)

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

func (r *RateLimiter) Close() {
	r.cancelCh <- struct{}{}
}
