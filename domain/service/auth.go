package service

import (
	"net"

	"github.com/teploff/antibruteforce/domain/entity"
)

type AuthService interface {
	LogIn(credentials entity.Credentials, ip net.IP) (bool, error)
}
