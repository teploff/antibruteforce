package grpc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/teploff/antibruteforce/internal/endpoints/auth"
	"github.com/teploff/antibruteforce/internal/transport/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxReceivedMsgSize = 1024 * 1024 * 20

var dur metrics.Histogram = kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
	Namespace: "antibruteforce",
	Subsystem: "auth",
	Name:      "latency_sign_in_request",
	Help:      "Total time (in seconds) spent serving requests by sign_in.",
}, []string{})

type server struct {
	signIn kitgrpc.Handler
}

func (s server) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInResponse, error) {
	defer func(begin time.Time) { dur.Observe(time.Since(begin).Seconds()) }(time.Now())

	_, response, err := s.signIn.ServeGRPC(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response.(*pb.SignInResponse), nil
}

// NewGRPCServer instance of gRPC server.
func NewGRPCServer(endpoints auth.Endpoints, errLogger log.Logger) *grpc.Server {
	options := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(errLogger)),
	}

	srv := &server{
		signIn: newRecoveryGRPCHandler(kitgrpc.NewServer(
			endpoints.SignIn,
			decodeSignInRequest,
			encodeSignInResponse,
			options...,
		), errLogger),
	}

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor), grpc.MaxRecvMsgSize(maxReceivedMsgSize))
	pb.RegisterAuthServer(baseServer, srv)

	return baseServer
}

func decodeSignInRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	request := grpcReq.(*pb.SignInRequest)

	return auth.SignInRequest{
		Login:    request.Login,
		Password: request.Password,
		IP:       request.Ip,
	}, nil
}

func encodeSignInResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(auth.SignInResponse)

	return &pb.SignInResponse{
		Ok: response.Ok,
	}, nil
}

//recoveryGRPCHandler wrap gRPC server, recover them if panic was fired.
type recoveryGRPCHandler struct {
	next   kitgrpc.Handler
	logger log.Logger
}

func newRecoveryGRPCHandler(next kitgrpc.Handler, logger log.Logger) *recoveryGRPCHandler {
	return &recoveryGRPCHandler{next: next, logger: logger}
}

func (rh *recoveryGRPCHandler) ServeGRPC(ctx context.Context, req interface{}) (context.Context, interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				_ = rh.logger.Log("msg", "gRPC server panic recover", "text", err.Error())
			}
		}
	}()

	return rh.next.ServeGRPC(ctx, req)
}
