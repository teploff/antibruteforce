package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type HTTPClient struct {
	Addr string
}

func NewHTTPClient(addr string) HTTPClient {
	return HTTPClient{Addr: addr}
}

func (c HTTPClient) ResetBucketByLogin(login string) error {
	request, err := encodeResetBucketByLoginRequest(login)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/reset_bucket_by_login", c.Addr)
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c HTTPClient) ResetBucketByPassword(password string) error {
	request, err := encodeResetBucketByPasswordRequest(password)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/reset_bucket_by_password", c.Addr)
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c HTTPClient) ResetBucketByIP(ip string) error {
	request, err := encodeResetBucketByIPRequest(ip)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/reset_bucket_by_ip", c.Addr)
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c HTTPClient) AddInBlacklist(subnet string) error {
	request, err := encodeSubnetRequest(subnet)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/add_in_blacklist", c.Addr)
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c HTTPClient) RemoveFromBlacklist(subnet string) error {
	request, err := encodeSubnetRequest(subnet)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/remove_from_blacklist", c.Addr)
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c HTTPClient) AddInWhitelist(subnet string) error {
	request, err := encodeSubnetRequest(subnet)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/add_in_whitelist", c.Addr)
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c HTTPClient) RemoveFromWhitelist(subnet string) error {
	request, err := encodeSubnetRequest(subnet)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/remove_from_whitelist", c.Addr)
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}
