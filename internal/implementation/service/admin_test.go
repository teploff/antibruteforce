package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/internal/implementation/repository/bucket"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
)

func TestResetBucketByLogin(t *testing.T) {
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
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	adminSvc := NewAdminService(ipList, loginBuckets, passwordBuckets, ipBuckets)
	_, err = loginBuckets.Add(credentials.Login)
	assert.NoError(t, err)
	assert.NoError(t, adminSvc.ResetBucketByLogin(credentials.Login))
	assert.Error(t, adminSvc.ResetBucketByLogin(credentials.Login))
	assert.NoError(t, flushAll(cfgRedis))
}

func TestResetBucketByPassword(t *testing.T) {
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
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	adminSvc := NewAdminService(ipList, loginBuckets, passwordBuckets, ipBuckets)
	_, err = passwordBuckets.Add(credentials.Password)
	assert.NoError(t, err)
	assert.NoError(t, adminSvc.ResetBucketByPassword(credentials.Password))
	assert.Error(t, adminSvc.ResetBucketByPassword(credentials.Password))
	assert.NoError(t, flushAll(cfgRedis))
}

func TestResetBucketByIP(t *testing.T) {
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
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	adminSvc := NewAdminService(ipList, loginBuckets, passwordBuckets, ipBuckets)
	_, err = ipBuckets.Add(IP.String())
	assert.NoError(t, err)
	assert.NoError(t, adminSvc.ResetBucketByIP(IP))
	assert.Error(t, adminSvc.ResetBucketByIP(IP))
	assert.NoError(t, flushAll(cfgRedis))
}

func TestAddInBlacklist(t *testing.T) {
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
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	adminSvc := NewAdminService(ipList, loginBuckets, passwordBuckets, ipBuckets)
	assert.NoError(t, adminSvc.AddInBlacklist(subNet))
	assert.Error(t, adminSvc.AddInBlacklist(subNet))
	assert.NoError(t, flushAll(cfgRedis))
}

func TestRemoveFromBlacklist(t *testing.T) {
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
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	adminSvc := NewAdminService(ipList, loginBuckets, passwordBuckets, ipBuckets)
	assert.NoError(t, adminSvc.AddInBlacklist(subNet))
	assert.NoError(t, adminSvc.RemoveFromBlacklist(subNet))
	assert.Error(t, adminSvc.RemoveFromBlacklist(subNet))
	assert.NoError(t, flushAll(cfgRedis))
}

func TestAddInWhitelist(t *testing.T) {
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
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	adminSvc := NewAdminService(ipList, loginBuckets, passwordBuckets, ipBuckets)
	assert.NoError(t, adminSvc.AddInWhitelist(subNet))
	assert.Error(t, adminSvc.AddInWhitelist(subNet))
	assert.NoError(t, flushAll(cfgRedis))
}

func TestRemoveFromWhitelist(t *testing.T) {
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
	ipList, err := ip.NewRedisIPList(cfgRedis)

	assert.NoError(t, err)
	assert.NoError(t, flushAll(cfgRedis))

	adminSvc := NewAdminService(ipList, loginBuckets, passwordBuckets, ipBuckets)
	assert.NoError(t, adminSvc.AddInWhitelist(subNet))
	assert.NoError(t, adminSvc.RemoveFromWhitelist(subNet))
	assert.Error(t, adminSvc.RemoveFromWhitelist(subNet))
	assert.NoError(t, flushAll(cfgRedis))
}
