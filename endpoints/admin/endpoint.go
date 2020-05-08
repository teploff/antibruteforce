package admin

import (
	"context"
	"net"

	"github.com/go-kit/kit/endpoint"
	"github.com/teploff/antibruteforce/domain/service"
)

// Endpoints for admin-panel.
type Endpoints struct {
	ResetBucketByLogin    endpoint.Endpoint
	ResetBucketByPassword endpoint.Endpoint
	ResetBucketByIP       endpoint.Endpoint
	AddInBlacklist        endpoint.Endpoint
	RemoveFromBlacklist   endpoint.Endpoint
	AddInWhitelist        endpoint.Endpoint
	RemoveFromWhitelist   endpoint.Endpoint
}

func makeResetBucketByLoginEndpoint(svc service.AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ResetBucketByLoginRequest)

		err = svc.ResetBucketByLogin(req.Login)
		if err != nil {
			return EmptyResponse{}, err
		}

		return EmptyResponse{}, nil
	}
}

func makeResetBucketByPasswordEndpoint(svc service.AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ResetBucketByPasswordRequest)

		err = svc.ResetBucketByPassword(req.Password)
		if err != nil {
			return EmptyResponse{}, err
		}

		return EmptyResponse{}, nil
	}
}

func makeResetBucketByIPEndpoint(svc service.AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ResetBucketByIPRequest)

		err = svc.ResetBucketByIP(net.ParseIP(req.IP))
		if err != nil {
			return EmptyResponse{}, err
		}

		return EmptyResponse{}, nil
	}
}

func makeAddInBlacklistEndpoint(svc service.AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SubnetRequest)

		_, ipNet, err := net.ParseCIDR(req.IPWithMask)
		if err != nil {
			return EmptyResponse{}, err
		}

		err = svc.AddInBlacklist(ipNet)
		if err != nil {
			return EmptyResponse{}, err
		}

		return EmptyResponse{}, nil
	}
}

func makeRemoveFromBlacklistEndpoint(svc service.AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SubnetRequest)

		_, ipNet, err := net.ParseCIDR(req.IPWithMask)
		if err != nil {
			return EmptyResponse{}, err
		}

		err = svc.RemoveFromBlacklist(ipNet)
		if err != nil {
			return EmptyResponse{}, err
		}

		return EmptyResponse{}, nil
	}
}

func makeAddInWhitelistEndpoint(svc service.AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SubnetRequest)

		_, ipNet, err := net.ParseCIDR(req.IPWithMask)
		if err != nil {
			return EmptyResponse{}, err
		}

		err = svc.AddInWhitelist(ipNet)
		if err != nil {
			return EmptyResponse{}, err
		}

		return EmptyResponse{}, nil
	}
}

func makeRemoveFromWhitelistEndpoint(svc service.AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SubnetRequest)

		_, ipNet, err := net.ParseCIDR(req.IPWithMask)
		if err != nil {
			return EmptyResponse{}, err
		}

		err = svc.RemoveFromWhitelist(ipNet)
		if err != nil {
			return EmptyResponse{}, err
		}

		return EmptyResponse{}, nil
	}
}

// MakeAdminEndpoints provides endpoints for admin-panel.
func MakeAdminEndpoints(svc service.AdminService) Endpoints {
	return Endpoints{
		ResetBucketByLogin:    makeResetBucketByLoginEndpoint(svc),
		ResetBucketByPassword: makeResetBucketByPasswordEndpoint(svc),
		ResetBucketByIP:       makeResetBucketByIPEndpoint(svc),
		AddInBlacklist:        makeAddInBlacklistEndpoint(svc),
		RemoveFromBlacklist:   makeRemoveFromBlacklistEndpoint(svc),
		AddInWhitelist:        makeAddInWhitelistEndpoint(svc),
		RemoveFromWhitelist:   makeRemoveFromWhitelistEndpoint(svc),
	}
}
