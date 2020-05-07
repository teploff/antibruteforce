package http

import (
	"bytes"
	"fmt"
	"net/http"
)

type Client struct {
	Addr string
}

func NewClient(addr string) Client {
	return Client{Addr: addr}
}

func (c Client) ResetBucketByLogin(login string) error {
	request, err := encodeResetBucketByLoginRequest(login)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/admin/reset_bucket_by_login", c.Addr)
	//nolint:gosec
	response, err := http.Post(url, "application/json", bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c Client) ResetBucketByPassword(password string) error {
	request, err := encodeResetBucketByPasswordRequest(password)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/admin/reset_bucket_by_password", c.Addr)
	//nolint:gosec
	response, err := http.Post(url, "application/json", bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c Client) ResetBucketByIP(ip string) error {
	request, err := encodeResetBucketByIPRequest(ip)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/admin/reset_bucket_by_ip", c.Addr)
	//nolint:gosec
	response, err := http.Post(url, "application/json", bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c Client) AddInBlacklist(subnet string) error {
	request, err := encodeSubnetRequest(subnet)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/add_in_blacklist", c.Addr)
	//nolint:gosec
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c Client) RemoveFromBlacklist(subnet string) error {
	request, err := encodeSubnetRequest(subnet)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/remove_from_blacklist", c.Addr)
	//nolint:gosec
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c Client) AddInWhitelist(subnet string) error {
	request, err := encodeSubnetRequest(subnet)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/add_in_whitelist", c.Addr)
	//nolint:gosec
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}

func (c Client) RemoveFromWhitelist(subnet string) error {
	request, err := encodeSubnetRequest(subnet)
	if err != nil {
		return err
	}

	destAddr := fmt.Sprintf("http://%s/admin/remove_from_whitelist", c.Addr)
	//nolint:gosec
	response, err := http.Post(destAddr, "application/json", bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return decodeResponse(response)
}
