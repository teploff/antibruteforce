package repository

import (
	"net"
)

// IPStorable provides storable interface for ip list for admin
//
// AddInWhitelist - adding subnet in the whitelist
//
// RemoveFromWhitelist - removing subnet from the whitelist
//
// AddInBlacklist - adding subnet in the blacklist
//
// RemoveFromBlacklist - removing subnet from the blacklist
//
// IsIPInWhiteList - checking, is ip in whitelist
//
// IsIPInBlackList - checking, is ip in blacklist
//
// WhiteListLength - returns length of the whitelist
//
// BlackListLength - returns length of the blacklist.
type IPStorable interface {
	AddInWhitelist(ipNet *net.IPNet) error
	RemoveFromWhitelist(ipNet *net.IPNet) error
	AddInBlacklist(ipNet *net.IPNet) error
	RemoveFromBlacklist(ipNet *net.IPNet) error
	IsIPInWhiteList(ip net.IP) (bool, error)
	IsIPInBlackList(ip net.IP) (bool, error)
	WhiteListLength() (int, error)
	BlackListLength() (int, error)
	Close() error
}
