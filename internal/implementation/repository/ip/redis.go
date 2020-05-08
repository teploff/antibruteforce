package ip

import (
	"errors"
	"net"

	"github.com/go-redis/redis/v7"
	pkgerrors "github.com/pkg/errors"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/domain/repository"
	"github.com/teploff/antibruteforce/internal/shared"
)

type redisIPList struct {
	whiteList *redis.Client
	blackList *redis.Client
}

// NewRedisIPList returns in redis repository of ip list.
func NewRedisIPList(cfg config.RedisConfig) (repository.IPStorable, error) {
	wl := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DbWhitelist,
	})

	if _, err := wl.Ping().Result(); err != nil {
		return nil, err
	}

	bl := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DbBlacklist,
	})

	if _, err := bl.Ping().Result(); err != nil {
		return nil, err
	}

	return &redisIPList{
		whiteList: wl,
		blackList: bl,
	}, nil
}

func (r *redisIPList) AddInWhitelist(ipNet *net.IPNet) error {
	count, err := r.blackList.Exists(ipNet.String()).Result()
	if err != nil {
		return err
	} else if count != 0 {
		return pkgerrors.Wrap(shared.ErrAlreadyExist, "in blacklist")
	}

	count, err = r.whiteList.Exists(ipNet.String()).Result()
	if err != nil {
		return err
	} else if count != 0 {
		return pkgerrors.Wrap(shared.ErrAlreadyExist, "in whitelist")
	}

	return r.whiteList.Set(ipNet.String(), ipNet.String(), 0).Err()
}

func (r *redisIPList) RemoveFromWhitelist(ipNet *net.IPNet) error {
	count, err := r.whiteList.Exists(ipNet.String()).Result()
	if err != nil {
		return err
	} else if count == 0 {
		return pkgerrors.Wrap(shared.ErrNotFound, "in whitelist")
	}

	err = r.whiteList.Del(ipNet.String()).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}

func (r *redisIPList) AddInBlacklist(ipNet *net.IPNet) error {
	count, err := r.whiteList.Exists(ipNet.String()).Result()
	if err != nil {
		return err
	} else if count != 0 {
		return pkgerrors.Wrap(shared.ErrAlreadyExist, "in whitelist")
	}

	count, err = r.blackList.Exists(ipNet.String()).Result()
	if err != nil {
		return err
	} else if count != 0 {
		return pkgerrors.Wrap(shared.ErrAlreadyExist, "in blacklist")
	}

	return r.blackList.Set(ipNet.String(), ipNet.String(), 0).Err()
}

func (r *redisIPList) RemoveFromBlacklist(ipNet *net.IPNet) error {
	count, err := r.blackList.Exists(ipNet.String()).Result()
	if err != nil {
		return err
	} else if count == 0 {
		return pkgerrors.Wrap(shared.ErrNotFound, "in blacklist")
	}

	err = r.blackList.Del(ipNet.String()).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}

func (r *redisIPList) IsIPInWhiteList(ip net.IP) (bool, error) {
	return r.isIPContains(ip, r.whiteList)
}

func (r *redisIPList) IsIPInBlackList(ip net.IP) (bool, error) {
	return r.isIPContains(ip, r.blackList)
}

func (r *redisIPList) WhiteListLength() (int, error) {
	res, err := r.whiteList.DBSize().Result()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *redisIPList) BlackListLength() (int, error) {
	res, err := r.blackList.DBSize().Result()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *redisIPList) isIPContains(ip net.IP, client *redis.Client) (bool, error) {
	keys, err := client.Keys("*").Result()
	if err != nil {
		return false, err
	}

	for _, key := range keys {
		_, ipNet, err := net.ParseCIDR(key)
		if err != nil {
			return false, err
		}

		if ipNet.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
}
