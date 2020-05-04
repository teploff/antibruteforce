package service

import (
	"fmt"
	"github.com/teploff/antibruteforce/domain/entity"
)

type authService struct{}

func NewAuthService() *authService {
	return &authService{}
}

func (a authService) LogIn(credentials entity.Credentials, ip string) (bool, error) {
	fmt.Printf("Login = %s; Password = %s; Ip = %s\n", credentials.Login, credentials.Password, ip)
	return true, nil
}
