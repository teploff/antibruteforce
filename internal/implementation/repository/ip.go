package repository

import (
	"net"
	"sync"

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

func (i *ipList) IsIPInWhiteList(ip net.IP) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.isIPContains(ip, i.whiteList)
}

func (i *ipList) IsIPInBlackList(ip net.IP) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.isIPContains(ip, i.blackList)
}

func (i *ipList) WhiteListLength() int {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return len(i.whiteList)
}

func (i *ipList) BlackListLength() int {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return len(i.blackList)
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
	for index := range list {
		if list[index].Contains(ip) {
			return true
		}
	}

	return false
}
