package http

import (
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/teploff/antibruteforce/internal/endpoints/admin"
	"go.uber.org/zap"
)

var (
	errRecovery = "error while handling request"
)

// NewHTTPServer instance of HTTP Handler.
func NewHTTPServer(endpoints admin.Endpoints, logger *zap.Logger) http.Handler {
	var r = mux.NewRouter()

	r.Methods("POST").Path("/admin/reset_bucket_by_login").Handler(withRecovery(kithttp.NewServer(
		endpoints.ResetBucketByLogin,
		decodeResetBucketByLoginRequest,
		encodeResponse,
	), logger))
	r.Methods("POST").Path("/admin/reset_bucket_by_password").Handler(withRecovery(kithttp.NewServer(
		endpoints.ResetBucketByPassword,
		decodeResetBucketByPasswordRequest,
		encodeResponse,
	), logger))
	r.Methods("POST").Path("/admin/reset_bucket_by_ip").Handler(withRecovery(kithttp.NewServer(
		endpoints.ResetBucketByIP,
		decodeResetBucketByIPRequest,
		encodeResponse,
	), logger))
	r.Methods("POST").Path("/admin/add_in_blacklist").Handler(withRecovery(kithttp.NewServer(
		endpoints.AddInBlacklist,
		decodeAddInBlacklistRequest,
		encodeResponse,
	), logger))
	r.Methods("POST").Path("/admin/remove_from_blacklist").Handler(withRecovery(kithttp.NewServer(
		endpoints.RemoveFromBlacklist,
		decodeRemoveFromBlacklistRequest,
		encodeResponse,
	), logger))
	r.Methods("POST").Path("/admin/add_in_whitelist").Handler(withRecovery(kithttp.NewServer(
		endpoints.AddInWhitelist,
		decodeAddInWhitelistRequest,
		encodeResponse,
	), logger))
	r.Methods("POST").Path("/admin/remove_from_whitelist").Handler(withRecovery(kithttp.NewServer(
		endpoints.RemoveFromWhitelist,
		decodeRemoveFromWhitelistRequest,
		encodeResponse,
	), logger))
	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	return r
}

func decodeResetBucketByLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request admin.ResetBucketByLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeResetBucketByPasswordRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request admin.ResetBucketByPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeResetBucketByIPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request admin.ResetBucketByIPRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeAddInBlacklistRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request admin.SubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeRemoveFromBlacklistRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request admin.SubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeAddInWhitelistRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request admin.SubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeRemoveFromWhitelistRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request admin.SubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// withRecovery handler middleware.
func withRecovery(h http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(errRecovery, zap.Any("error", err))
				http.Error(w, errRecovery, http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
