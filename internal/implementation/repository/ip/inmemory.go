package ip

import (
	"net"
	"sync"

	"github.com/pkg/errors"
	"github.com/teploff/antibruteforce/domain/repository"
	"github.com/teploff/antibruteforce/internal/shared"
)

type ipList struct {
	whiteList []*net.IPNet
	blackList []*net.IPNet
	mu        *sync.RWMutex
}

func NewIPList() repository.IPStorable {
	return &ipList{
		whiteList: make([]*net.IPNet, 0, 10),
		blackList: make([]*net.IPNet, 0, 10),
		mu:        &sync.RWMutex{},
	}
}

func (i *ipList) AddInWhitelist(ipNet *net.IPNet) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exist := i.findElem(ipNet, i.whiteList); exist {
		return shared.ErrAlreadyExist
	}

	if _, exist := i.findElem(ipNet, i.blackList); exist {
		return errors.Wrap(shared.ErrAlreadyExist, "in blacklist")
	}

	i.whiteList = append(i.whiteList, ipNet)

	return nil
}

func (i *ipList) RemoveFromWhitelist(ipNet *net.IPNet) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	index, exist := i.findElem(ipNet, i.whiteList)
	if !exist {
		return shared.ErrNotFound
	}

	i.whiteList = append(i.whiteList[:index], i.whiteList[index+1:]...)

	return nil
}

func (i *ipList) AddInBlacklist(ipNet *net.IPNet) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exist := i.findElem(ipNet, i.blackList); exist {
		return shared.ErrAlreadyExist
	}

	if _, exist := i.findElem(ipNet, i.whiteList); exist {
		return errors.Wrap(shared.ErrAlreadyExist, "in whitelist")
	}

	i.blackList = append(i.blackList, ipNet)

	return nil
}

func (i *ipList) RemoveFromBlacklist(ipNet *net.IPNet) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	index, exist := i.findElem(ipNet, i.blackList)
	if !exist {
		return shared.ErrNotFound
	}

	i.blackList = append(i.blackList[:index], i.blackList[index+1:]...)

	return nil
}

func (i *ipList) IsIPInWhiteList(ip net.IP) (bool, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.isIPContains(ip, i.whiteList), nil
}

func (i *ipList) IsIPInBlackList(ip net.IP) (bool, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.isIPContains(ip, i.blackList), nil
}

func (i *ipList) WhiteListLength() (int, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return len(i.whiteList), nil
}

func (i *ipList) BlackListLength() (int, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return len(i.blackList), nil
}

func (i *ipList) findElem(ipNet *net.IPNet, list []*net.IPNet) (int, bool) {
	for index := range list {
		if list[index].IP.String() == ipNet.IP.String() && list[index].Mask.String() == ipNet.Mask.String() {
			return index, true
		}
	}

	return -1, false
}

func (i *ipList) isIPContains(ip net.IP, list []*net.IPNet) bool {
	for _, ipNet := range list {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}
