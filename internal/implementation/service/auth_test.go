package service

import (
	"net"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/internal/implementation/repository/bucket"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
	"github.com/teploff/antibruteforce/internal/limiter"
)

var (
	cfgRedis = config.RedisConfig{
		Addr:        "0.0.0.0:6379",
		Password:    "",
		DbWhitelist: 14,
		DbBlacklist: 15,
	}
	credentials = entity.Credentials{
		Login:    "login",
		Password: "password",
	}
	IP           = net.ParseIP("192.168.130.132")
	_, subNet, _ = net.ParseCIDR("192.168.130.0/24")
)

func TestIPInWhitelist(t *testing.T) {
	cfgRL := config.RateLimiterConfig{
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

	loginBuckets := bucket.NewLeakyBucket(cfgRL.Login.Rate, cfgRL.Login.Interval, cfgRL.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfgRL.Password.Rate, cfgRL.Password.Interval, cfgRL.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfgRL.IP.Rate, cfgRL.IP.Interval, cfgRL.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfgRL.GCTime)
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	authSvc := NewAuthService(rateLimiter, ipList)
	assert.NoError(t, ipList.AddInWhitelist(subNet))

	allow, err := authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
}

func TestIPInBlacklist(t *testing.T) {
	cfgRL := config.RateLimiterConfig{
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

	loginBuckets := bucket.NewLeakyBucket(cfgRL.Login.Rate, cfgRL.Login.Interval, cfgRL.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfgRL.Password.Rate, cfgRL.Password.Interval, cfgRL.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfgRL.IP.Rate, cfgRL.IP.Interval, cfgRL.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfgRL.GCTime)
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	authSvc := NewAuthService(rateLimiter, ipList)
	assert.NoError(t, ipList.AddInBlacklist(subNet))

	allow, err := authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.False(t, allow)

	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.False(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.False(t, allow)
}

func TestBruteForceByLogin(t *testing.T) {
	cfgRL := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       2,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       10,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		IP: config.IP{
			Rate:       10,
			Interval:   time.Second,
			ExpireTime: time.Minute * 1,
		},
		GCTime: time.Minute * 20,
	}

	loginBuckets := bucket.NewLeakyBucket(cfgRL.Login.Rate, cfgRL.Login.Interval, cfgRL.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfgRL.Password.Rate, cfgRL.Password.Interval, cfgRL.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfgRL.IP.Rate, cfgRL.IP.Interval, cfgRL.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfgRL.GCTime)
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	authSvc := NewAuthService(rateLimiter, ipList)
	allow, err := authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.False(t, allow)
}

func TestBruteForceByPassword(t *testing.T) {
	cfgRL := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       10,
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
			Interval:   time.Second,
			ExpireTime: time.Minute * 1,
		},
		GCTime: time.Minute * 20,
	}

	loginBuckets := bucket.NewLeakyBucket(cfgRL.Login.Rate, cfgRL.Login.Interval, cfgRL.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfgRL.Password.Rate, cfgRL.Password.Interval, cfgRL.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfgRL.IP.Rate, cfgRL.IP.Interval, cfgRL.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfgRL.GCTime)
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	authSvc := NewAuthService(rateLimiter, ipList)
	allow, err := authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.False(t, allow)
}

func TestBruteForceByIP(t *testing.T) {
	cfgRL := config.RateLimiterConfig{
		Login: config.Login{
			Rate:       10,
			Interval:   time.Second,
			ExpireTime: time.Minute * 10,
		},
		Password: config.Password{
			Rate:       10,
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

	loginBuckets := bucket.NewLeakyBucket(cfgRL.Login.Rate, cfgRL.Login.Interval, cfgRL.Login.ExpireTime)
	passwordBuckets := bucket.NewLeakyBucket(cfgRL.Password.Rate, cfgRL.Password.Interval, cfgRL.Password.ExpireTime)
	ipBuckets := bucket.NewLeakyBucket(cfgRL.IP.Rate, cfgRL.IP.Interval, cfgRL.IP.ExpireTime)
	rateLimiter := limiter.NewRateLimiter(loginBuckets, passwordBuckets, ipBuckets, cfgRL.GCTime)
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	authSvc := NewAuthService(rateLimiter, ipList)
	allow, err := authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.True(t, allow)
	allow, err = authSvc.LogIn(credentials, IP)
	assert.NoError(t, err)
	assert.False(t, allow)
}

func flushAll(cfg config.RedisConfig) error {
	wl := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DbWhitelist,
	})

	if _, err := wl.Ping().Result(); err != nil {
		return err
	}

	_, err := wl.FlushAll().Result()

	if err != nil {
		return err
	}

	return nil
}
