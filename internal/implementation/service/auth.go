package service

import (
	"errors"

	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/gateway"
	"github.com/teploff/antibruteforce/domain/service"
	"github.com/teploff/antibruteforce/internal/shared"
)

type authService struct {
	loginRl    gateway.RateLimiter
	passwordRl gateway.RateLimiter
	ipRl       gateway.RateLimiter
}

func NewAuthService(loginRl, passwordRl, ipRl gateway.RateLimiter) service.AuthService {
	return &authService{
		loginRl:    loginRl,
		passwordRl: passwordRl,
		ipRl:       ipRl,
	}
}

func (a *authService) LogIn(credentials entity.Credentials, ip string) (bool, error) {
	bruteForce, err := isBruteForce(a.loginRl, credentials.Login)
	if err != nil {
		return false, err
	}

	if bruteForce {
		return false, nil
	}

	bruteForce, err = isBruteForce(a.passwordRl, credentials.Password)
	if err != nil {
		return false, err
	}

	if bruteForce {
		return false, nil
	}

	bruteForce, err = isBruteForce(a.ipRl, ip)
	if err != nil {
		return false, err
	}

	if bruteForce {
		return false, nil
	}

	return true, nil
}

func isBruteForce(rate gateway.RateLimiter, keyBucket string) (bool, error) {
	limiter, err := rate.GetLimiter(keyBucket)
	if errors.Is(err, shared.ErrNotFound) {
		if limiter, err = rate.AddBucket(keyBucket); err != nil {
			return false, err
		}
	}

	if !limiter.Allow() {
		return true, nil
	}

	return false, nil
}
