package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/service"
)

type Endpoints struct {
	SignIn endpoint.Endpoint
}

func makeSignInEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SignInRequest)
		ok, err := svc.LogIn(entity.Credentials{
			Login:    req.Login,
			Password: req.Password,
		}, req.IP)

		if err != nil {
			return SignInResponse{}, err
		}

		return SignInResponse{
			Ok: ok,
		}, nil
	}
}

func MakeAuthEndpoints(svc service.AuthService) Endpoints {
	return Endpoints{
		SignIn: makeSignInEndpoint(svc),
	}
}
