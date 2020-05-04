package service

import "github.com/teploff/antibruteforce/domain/entity"

type AuthService interface {
	LogIn(credentials entity.Credentials, ip string) (bool, error)
}
