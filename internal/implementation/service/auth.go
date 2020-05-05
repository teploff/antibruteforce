package service

import (
	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/service"
	"github.com/teploff/antibruteforce/internal/limiter"
)

type authService struct {
	rl *limiter.RateLimiter
}

func NewAuthService(rateLimiter *limiter.RateLimiter) service.AuthService {
	return &authService{rl: rateLimiter}
}

func (a *authService) LogIn(credentials entity.Credentials, ip string) (bool, error) {
	bruteForce, err := a.rl.IsBruteForce(credentials.Login, credentials.Password, ip)
	if err != nil || bruteForce {
		return false, err
	}

	return true, nil
}
