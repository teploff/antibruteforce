package service

import (
	"fmt"

	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/service"
)

type authService struct{}

func NewAuthService() service.AuthService {
	return &authService{}
}

func (a authService) LogIn(credentials entity.Credentials, ip string) (bool, error) {
	fmt.Printf("Login = %s; Password = %s; Ip = %s\n", credentials.Login, credentials.Password, ip)
	return true, nil
}
