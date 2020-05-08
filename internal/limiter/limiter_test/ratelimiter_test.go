package limiter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/internal/implementation/repository/bucket"
	"github.com/teploff/antibruteforce/internal/limiter"
)

const (
	login    = "login"
	password = "password"
	ip       = "192.168.130.132"
)

func TestBruteForceOnOnlyLoginBuckets(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		IP: config.IP{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		GCTime: time.Minute * 20,
	}
	loginBuckets := bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.GCTime)

	isBruteForce, err := rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.True(t, isBruteForce)
	assert.NoError(t, err)
}

func TestBruteForceOnOnlyPasswordBuckets(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		IP: config.IP{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		GCTime: time.Minute * 20,
	}
	loginBuckets := bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.GCTime)

	isBruteForce, err := rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.True(t, isBruteForce)
	assert.NoError(t, err)
}

func TestBruteForceOnOnlyIpBuckets(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		IP: config.IP{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		GCTime: time.Minute * 20,
	}
	loginBuckets := bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.GCTime)

	isBruteForce, err := rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.True(t, isBruteForce)
	assert.NoError(t, err)
}

func TestBruteForceOnOnlyLoginAndPasswordBuckets(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		IP: config.IP{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		GCTime: time.Minute * 20,
	}
	loginBuckets := bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.GCTime)

	isBruteForce, err := rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.True(t, isBruteForce)
	assert.NoError(t, err)
}

func TestBruteForceOnOnlyLoginAndIPBuckets(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		IP: config.IP{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		GCTime: time.Minute * 20,
	}
	loginBuckets := bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.GCTime)

	isBruteForce, err := rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.True(t, isBruteForce)
	assert.NoError(t, err)
}

func TestBruteForceOnOnlyPasswordAndIPBuckets(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       10,
			Interval:   time.Minute,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		IP: config.IP{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		GCTime: time.Minute * 20,
	}
	loginBuckets := bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.GCTime)

	isBruteForce, err := rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.True(t, isBruteForce)
	assert.NoError(t, err)
}

func TestBruteForceDoesntExist(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		IP: config.IP{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 1,
		},
		GCTime: time.Minute * 20,
	}
	loginBuckets := bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.GCTime)

	isBruteForce, err := rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)

	time.Sleep(time.Second)

	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.True(t, isBruteForce)
	assert.NoError(t, err)
}

func TestExpiredBucket(t *testing.T) {
	cfg := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Millisecond * 500,
		},
		Password: config.Password{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Millisecond * 500,
		},
		IP: config.IP{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Millisecond * 500,
		},
		GCTime: time.Millisecond * 100,
	}
	loginBuckets := bucket.NewLeakyBucket(cfg.Login.Rate, cfg.Login.Interval, cfg.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfg.Password.Rate, cfg.Password.Interval, cfg.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfg.IP.Rate, cfg.IP.Interval, cfg.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfg.GCTime)
	go rateLimiter.RunGarbageCollector()

	isBruteForce, err := rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 700)

	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.False(t, isBruteForce)
	assert.NoError(t, err)
	isBruteForce, err = rateLimiter.IsBruteForce(login, password, ip)
	assert.True(t, isBruteForce)
	assert.NoError(t, err)

	rateLimiter.Close()
}
