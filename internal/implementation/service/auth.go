package service

import (
	"net"

	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/repository"
	"github.com/teploff/antibruteforce/domain/service"
	"github.com/teploff/antibruteforce/internal/limiter"
)

type authService struct {
	rl     *limiter.RateLimiter
	ipList repository.IPStorable
}

func NewAuthService(rateLimiter *limiter.RateLimiter, ipList repository.IPStorable) service.AuthService {
	return &authService{
		rl:     rateLimiter,
		ipList: ipList,
	}
}

func (a *authService) LogIn(credentials entity.Credentials, ip net.IP) (bool, error) {
	inList, err := a.ipList.IsIPInWhiteList(ip)
	if err != nil {
		return false, err
	}

	if inList {
		return true, nil
	}

	inList, err = a.ipList.IsIPInBlackList(ip)
	if err != nil {
		return false, err
	}

	if inList {
		return false, nil
	}

	bruteForce, err := a.rl.IsBruteForce(credentials.Login, credentials.Password, ip.String())
	if err != nil || bruteForce {
		return false, err
	}

	return true, nil
}
