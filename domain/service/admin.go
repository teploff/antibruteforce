package service

import "net"

type AdminService interface {
	ResetBucketByLogin(login string) error
	ResetBucketByPassword(password string) error
	ResetBucketByIP(ip net.IP) error
	AddInBlacklist(ipNet *net.IPNet) error
	RemoveFromBlacklist(ipNet *net.IPNet) error
	AddInWhitelist(ipNet *net.IPNet) error
	RemoveFromWhitelist(ipNet *net.IPNet) error
}
