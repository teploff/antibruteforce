package service

import (
	"net"

	"github.com/teploff/antibruteforce/domain/repository"
	"github.com/teploff/antibruteforce/domain/service"
)

// adminService implementation of admin service.
type adminService struct {
	ipList          repository.IPStorable
	loginBuckets    repository.BucketStorable
	passwordBuckets repository.BucketStorable
	ipBuckets       repository.BucketStorable
}

// NewAdminService returns implementation of admin service.
func NewAdminService(ipList repository.IPStorable, login, password, ip repository.BucketStorable) service.AdminService {
	return &adminService{
		ipList:          ipList,
		loginBuckets:    login,
		passwordBuckets: password,
		ipBuckets:       ip,
	}
}

func (a adminService) ResetBucketByLogin(login string) error {
	return a.loginBuckets.Delete(login)
}

func (a adminService) ResetBucketByPassword(password string) error {
	return a.passwordBuckets.Delete(password)
}

func (a adminService) ResetBucketByIP(ip net.IP) error {
	return a.ipBuckets.Delete(ip.String())
}

func (a adminService) AddInBlacklist(ipNet *net.IPNet) error {
	return a.ipList.AddInBlacklist(ipNet)
}

func (a adminService) RemoveFromBlacklist(ipNet *net.IPNet) error {
	return a.ipList.RemoveFromBlacklist(ipNet)
}

func (a adminService) AddInWhitelist(ipNet *net.IPNet) error {
	return a.ipList.AddInWhitelist(ipNet)
}

func (a adminService) RemoveFromWhitelist(ipNet *net.IPNet) error {
	return a.ipList.RemoveFromWhitelist(ipNet)
}
