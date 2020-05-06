package repository_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teploff/antibruteforce/internal/implementation/repository"
)

//nolint:dupl
func TestIncrementWhiteAndBlackLists(t *testing.T) {
	list := repository.NewIPList()
	_, whiteNet1, _ := net.ParseCIDR("192.168.128.0/24")
	_, whiteNet2, _ := net.ParseCIDR("192.168.129.0/24")
	_, whiteNet3, _ := net.ParseCIDR("192.168.130.0/24")

	_, blackNet1, _ := net.ParseCIDR("192.168.131.0/24")
	_, blackNet2, _ := net.ParseCIDR("192.168.132.0/24")
	_, blackNet3, _ := net.ParseCIDR("192.168.133.0/24")

	assert.NoError(t, list.AddInWhitelist(whiteNet1))
	assert.NoError(t, list.AddInWhitelist(whiteNet2))
	assert.NoError(t, list.AddInWhitelist(whiteNet3))
	assert.Error(t, list.AddInWhitelist(whiteNet1))
	assert.Error(t, list.AddInWhitelist(whiteNet2))
	assert.Error(t, list.AddInWhitelist(whiteNet3))

	assert.NoError(t, list.AddInBlacklist(blackNet1))
	assert.NoError(t, list.AddInBlacklist(blackNet2))
	assert.NoError(t, list.AddInBlacklist(blackNet3))
	assert.Error(t, list.AddInBlacklist(blackNet1))
	assert.Error(t, list.AddInBlacklist(blackNet2))
	assert.Error(t, list.AddInBlacklist(blackNet3))

	assert.Equal(t, 3, list.WhiteListLength())
	assert.Equal(t, 3, list.BlackListLength())
}

//nolint:dupl
func TestRemovingWhiteAndBlackLists(t *testing.T) {
	list := repository.NewIPList()
	_, whiteNet1, _ := net.ParseCIDR("192.168.128.0/24")
	_, whiteNet2, _ := net.ParseCIDR("192.168.129.0/24")
	_, whiteNet3, _ := net.ParseCIDR("192.168.130.0/24")

	_, blackNet1, _ := net.ParseCIDR("192.168.131.0/24")
	_, blackNet2, _ := net.ParseCIDR("192.168.132.0/24")
	_, blackNet3, _ := net.ParseCIDR("192.168.133.0/24")

	assert.NoError(t, list.AddInWhitelist(whiteNet1))
	assert.NoError(t, list.AddInWhitelist(whiteNet2))
	assert.NoError(t, list.AddInWhitelist(whiteNet3))

	assert.NoError(t, list.AddInBlacklist(blackNet1))
	assert.NoError(t, list.AddInBlacklist(blackNet2))
	assert.NoError(t, list.AddInBlacklist(blackNet3))

	assert.NoError(t, list.RemoveFromWhitelist(whiteNet1))
	assert.NoError(t, list.RemoveFromWhitelist(whiteNet2))
	assert.NoError(t, list.RemoveFromWhitelist(whiteNet3))

	assert.NoError(t, list.RemoveFromBlacklist(blackNet1))
	assert.NoError(t, list.RemoveFromBlacklist(blackNet2))
	assert.NoError(t, list.RemoveFromBlacklist(blackNet3))

	assert.Equal(t, 0, list.WhiteListLength())
	assert.Equal(t, 0, list.BlackListLength())
}

//nolint:funlen
func TestBelongWhiteAndBlackLists(t *testing.T) {
	list := repository.NewIPList()
	_, whiteNet1, _ := net.ParseCIDR("192.168.130.0/24")
	_, whiteNet2, _ := net.ParseCIDR("192.168.0.0/16")
	_, whiteNet3, _ := net.ParseCIDR("192.0.0.0/8")

	_, blackNet1, _ := net.ParseCIDR("10.200.128.0/24")
	_, blackNet2, _ := net.ParseCIDR("10.200.0.0/16")
	_, blackNet3, _ := net.ParseCIDR("10.0.0.0/8")

	whiteIP1 := net.ParseIP("192.15.10.11")
	whiteIP2 := net.ParseIP("192.168.10.11")
	whiteIP3 := net.ParseIP("192.168.130.11")

	blackIP1 := net.ParseIP("10.15.10.11")
	blackIP2 := net.ParseIP("10.200.10.11")
	blackIP3 := net.ParseIP("10.200.128.11")

	neutralIP := net.ParseIP("127.0.0.1")

	assert.NoError(t, list.AddInWhitelist(whiteNet1))
	assert.False(t, list.IsIPInWhiteList(whiteIP1))
	assert.False(t, list.IsIPInWhiteList(whiteIP2))
	assert.True(t, list.IsIPInWhiteList(whiteIP3))
	assert.NoError(t, list.RemoveFromWhitelist(whiteNet1))

	assert.NoError(t, list.AddInWhitelist(whiteNet2))
	assert.False(t, list.IsIPInWhiteList(whiteIP1))
	assert.True(t, list.IsIPInWhiteList(whiteIP2))
	assert.True(t, list.IsIPInWhiteList(whiteIP3))
	assert.NoError(t, list.RemoveFromWhitelist(whiteNet2))

	assert.NoError(t, list.AddInWhitelist(whiteNet3))
	assert.True(t, list.IsIPInWhiteList(whiteIP1))
	assert.True(t, list.IsIPInWhiteList(whiteIP2))
	assert.True(t, list.IsIPInWhiteList(whiteIP3))
	assert.NoError(t, list.RemoveFromWhitelist(whiteNet3))

	assert.NoError(t, list.AddInBlacklist(blackNet1))
	assert.False(t, list.IsIPInBlackList(blackIP1))
	assert.False(t, list.IsIPInBlackList(blackIP2))
	assert.True(t, list.IsIPInBlackList(blackIP3))
	assert.NoError(t, list.RemoveFromBlacklist(blackNet1))

	assert.NoError(t, list.AddInBlacklist(blackNet2))
	assert.False(t, list.IsIPInBlackList(blackIP1))
	assert.True(t, list.IsIPInBlackList(blackIP2))
	assert.True(t, list.IsIPInBlackList(blackIP3))
	assert.NoError(t, list.RemoveFromBlacklist(blackNet2))

	assert.NoError(t, list.AddInBlacklist(blackNet3))
	assert.True(t, list.IsIPInBlackList(blackIP1))
	assert.True(t, list.IsIPInBlackList(blackIP2))
	assert.True(t, list.IsIPInBlackList(blackIP3))
	assert.NoError(t, list.RemoveFromBlacklist(blackNet3))

	assert.NoError(t, list.AddInWhitelist(whiteNet1))
	assert.NoError(t, list.AddInWhitelist(whiteNet2))
	assert.NoError(t, list.AddInWhitelist(whiteNet3))

	assert.NoError(t, list.AddInBlacklist(blackNet1))
	assert.NoError(t, list.AddInBlacklist(blackNet2))
	assert.NoError(t, list.AddInBlacklist(blackNet3))

	assert.False(t, list.IsIPInWhiteList(neutralIP))
	assert.False(t, list.IsIPInBlackList(neutralIP))

	assert.Equal(t, 3, list.WhiteListLength())
	assert.Equal(t, 3, list.BlackListLength())
}
