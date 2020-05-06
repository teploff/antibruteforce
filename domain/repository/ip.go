package repository

import (
	"net"
)

type IPStorable interface {
	AddInWhitelist(ipNet *net.IPNet) error
	RemoveFromWhitelist(ipNet *net.IPNet) error
	AddInBlacklist(ipNet *net.IPNet) error
	RemoveFromBlacklist(ipNet *net.IPNet) error
	IsIPInWhiteList(ip net.IP) bool
	IsIPInBlackList(ip net.IP) bool
	WhiteListLength() int
	BlackListLength() int
}
