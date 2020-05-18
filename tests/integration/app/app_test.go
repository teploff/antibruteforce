package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/repository"
	"github.com/teploff/antibruteforce/endpoints/admin"
	"github.com/teploff/antibruteforce/internal/app"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
	"github.com/teploff/antibruteforce/internal/shared"
	"github.com/teploff/antibruteforce/transport/grpc/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type GRPCBruteForceTestSuit struct {
	suite.Suite
	cfg         config.Config
	client      *mongo.Client
	ipList      repository.IPStorable
	credentials entity.Credentials
	userIP      net.IP
	application *app.App
}

func (g *GRPCBruteForceTestSuit) SetupSuite() {
	cfg, err := config.LoadFromFile("../../../init/config_test.yaml")
	if err != nil {
		panic(err)
	}

	g.cfg = cfg
	g.cfg.Mongo.DBName = "gRPC_brute_force"

	g.credentials = entity.Credentials{
		Login:    "loginAuth",
		Password: "passwordAuth",
	}

	g.userIP = net.ParseIP("192.168.199.132")

	g.client, _ = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", g.cfg.Mongo.Addr)))
	_ = g.client.Connect(context.TODO())

	g.ipList, _ = ip.NewMongoIPList(g.cfg.Mongo)

	g.application = app.NewApp(g.cfg,
		app.WithLeakyBuckets(g.cfg.RateLimiter),
		app.WithIPList(g.ipList),
	)

	go g.application.Run()
}

func (g *GRPCBruteForceTestSuit) TearDownSuite() {
	_ = g.ipList.Close()
	_ = g.client.Disconnect(context.TODO())
}

func (g *GRPCBruteForceTestSuit) TearDownTest() {
	_ = g.client.Database(g.cfg.Mongo.DBName).Collection("whitelist").Drop(context.TODO())
	_ = g.client.Database(g.cfg.Mongo.DBName).Collection("blacklist").Drop(context.TODO())

	time.Sleep(time.Millisecond * 50)
}

func TestGRPCBruteForce(t *testing.T) {
	suite.Run(t, new(GRPCBruteForceTestSuit))
}

func (g *GRPCBruteForceTestSuit) TestBruteForceByLogin() {
	conn, err := grpc.Dial(g.cfg.GRPCServer.Addr, grpc.WithInsecure())
	g.Assert().NoError(err)

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    g.credentials.Login,
		Password: g.credentials.Password,
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    g.credentials.Login,
		Password: g.credentials.Password,
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    g.credentials.Login,
		Password: g.credentials.Password,
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().False(response.Ok)
}

func (g *GRPCBruteForceTestSuit) TestBruteForceByPassword() {
	conn, err := grpc.Dial(g.cfg.GRPCServer.Addr, grpc.WithInsecure())
	g.Assert().NoError(err)

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login1",
		Password: g.credentials.Password,
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login2",
		Password: g.credentials.Password,
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login3",
		Password: g.credentials.Password,
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login4",
		Password: g.credentials.Password,
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().False(response.Ok)
}

func (g *GRPCBruteForceTestSuit) TestBruteForceByIP() {
	conn, err := grpc.Dial(g.cfg.GRPCServer.Addr, grpc.WithInsecure())
	g.Assert().NoError(err)

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login1",
		Password: "password1",
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login2",
		Password: "password2",
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login3",
		Password: "password3",
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login4",
		Password: "password4",
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login5",
		Password: "password5",
		Ip:       g.userIP.String(),
	})
	g.Assert().NoError(err)
	g.Assert().False(response.Ok)
}

type HTTPAdminPanelTestSuit struct {
	suite.Suite
	cfg         config.Config
	client      *mongo.Client
	ipList      repository.IPStorable
	credentials entity.Credentials
	userIP      net.IP
	subnet      *net.IPNet
	app         *app.App
}

func (h *HTTPAdminPanelTestSuit) SetupSuite() {
	cfg, err := config.LoadFromFile("../../../init/config_test.yaml")
	if err != nil {
		panic(err)
	}

	h.cfg = cfg
	h.cfg.Mongo.DBName = "gRPC_http_admin_panel"
	h.cfg.GRPCServer.Addr = "localhost:8114"
	h.cfg.HTTPServer.Addr = "localhost:8115"

	h.credentials = entity.Credentials{
		Login:    "loginAuth",
		Password: "passwordAuth",
	}

	h.userIP = net.ParseIP("192.168.199.132")
	_, h.subnet, _ = net.ParseCIDR("192.168.199.0/24")

	h.client, _ = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", h.cfg.Mongo.Addr)))
	_ = h.client.Connect(context.TODO())

	h.ipList, _ = ip.NewMongoIPList(h.cfg.Mongo)

	h.app = app.NewApp(h.cfg,
		app.WithLeakyBuckets(h.cfg.RateLimiter),
		app.WithIPList(h.ipList),
	)

	go h.app.Run()
}

func (h *HTTPAdminPanelTestSuit) TearDownSuite() {
	_ = h.ipList.Close()
	_ = h.client.Disconnect(context.TODO())
}

func (h *HTTPAdminPanelTestSuit) TearDownTest() {
	_ = h.client.Database(h.cfg.Mongo.DBName).Collection("whitelist").Drop(context.TODO())
	_ = h.client.Database(h.cfg.Mongo.DBName).Collection("blacklist").Drop(context.TODO())

	time.Sleep(time.Millisecond * 50)
}

func TestHTTPAdminPanel(t *testing.T) {
	suite.Run(t, new(HTTPAdminPanelTestSuit))
}

func (h *HTTPAdminPanelTestSuit) TestResetByLogin() {
	conn, err := grpc.Dial(h.cfg.GRPCServer.Addr, grpc.WithInsecure())
	h.Assert().NoError(err)

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: "password1",
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: "password2",
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: "password3",
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().False(response.Ok)

	marshalReq, err := json.Marshal(admin.ResetBucketByLoginRequest{Login: h.credentials.Login})
	h.Assert().NoError(err)
	url := fmt.Sprintf("http://%s/admin/reset_bucket_by_login", h.cfg.HTTPServer.Addr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	h.Assert().NoError(err)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().Error(decodeResponse(resp))

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)
}

func (h *HTTPAdminPanelTestSuit) TestResetByPassword() {
	conn, err := grpc.Dial(h.cfg.GRPCServer.Addr, grpc.WithInsecure())
	h.Assert().NoError(err)

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login1",
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login2",
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login3",
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login4",
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().False(response.Ok)

	marshalReq, err := json.Marshal(admin.ResetBucketByPasswordRequest{Password: h.credentials.Password})
	h.Assert().NoError(err)
	url := fmt.Sprintf("http://%s/admin/reset_bucket_by_password", h.cfg.HTTPServer.Addr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	h.Assert().NoError(err)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().Error(decodeResponse(resp))

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)
}

func (h *HTTPAdminPanelTestSuit) TestResetByIP() {
	conn, err := grpc.Dial(h.cfg.GRPCServer.Addr, grpc.WithInsecure())
	h.Assert().NoError(err)

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login1",
		Password: "password1",
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login2",
		Password: "password2",
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login3",
		Password: "password3",
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login4",
		Password: "password4",
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    "login5",
		Password: "password5",
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().False(response.Ok)

	marshalReq, err := json.Marshal(admin.ResetBucketByIPRequest{IP: h.userIP.String()})
	h.Assert().NoError(err)
	url := fmt.Sprintf("http://%s/admin/reset_bucket_by_ip", h.cfg.HTTPServer.Addr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	h.Assert().NoError(err)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().Error(decodeResponse(resp))

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)
}

func (h *HTTPAdminPanelTestSuit) TestAddInBlacklist() {
	conn, err := grpc.Dial(h.cfg.GRPCServer.Addr, grpc.WithInsecure())
	h.Assert().NoError(err)

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	marshalReq, err := json.Marshal(admin.SubnetRequest{IPWithMask: h.subnet.String()})
	h.Assert().NoError(err)
	url := fmt.Sprintf("http://%s/admin/add_in_blacklist", h.cfg.HTTPServer.Addr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	h.Assert().NoError(err)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().Error(decodeResponse(resp))

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().False(response.Ok)
}

func (h *HTTPAdminPanelTestSuit) TestRemoveFromBlacklist() {
	conn, err := grpc.Dial(h.cfg.GRPCServer.Addr, grpc.WithInsecure())
	h.Assert().NoError(err)

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	marshalReq, err := json.Marshal(admin.SubnetRequest{IPWithMask: h.subnet.String()})
	h.Assert().NoError(err)
	url := fmt.Sprintf("http://%s/admin/add_in_blacklist", h.cfg.HTTPServer.Addr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	h.Assert().NoError(err)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().Error(decodeResponse(resp))

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().False(response.Ok)

	marshalReq, err = json.Marshal(admin.SubnetRequest{IPWithMask: h.subnet.String()})
	h.Assert().NoError(err)
	url = fmt.Sprintf("http://%s/admin/remove_from_blacklist", h.cfg.HTTPServer.Addr)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)
}

func (h *HTTPAdminPanelTestSuit) TestAddInWhitelist() {
	conn, err := grpc.Dial(h.cfg.GRPCServer.Addr, grpc.WithInsecure())
	h.Assert().NoError(err)

	marshalReq, err := json.Marshal(admin.SubnetRequest{IPWithMask: h.subnet.String()})
	h.Assert().NoError(err)
	url := fmt.Sprintf("http://%s/admin/add_in_whitelist", h.cfg.HTTPServer.Addr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	h.Assert().NoError(err)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().Error(decodeResponse(resp))

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)
}

func (h *HTTPAdminPanelTestSuit) TestRemoveFromWhitelist() {
	conn, err := grpc.Dial(h.cfg.GRPCServer.Addr, grpc.WithInsecure())
	h.Assert().NoError(err)

	marshalReq, err := json.Marshal(admin.SubnetRequest{IPWithMask: h.subnet.String()})
	h.Assert().NoError(err)
	url := fmt.Sprintf("http://%s/admin/add_in_whitelist", h.cfg.HTTPServer.Addr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	h.Assert().NoError(err)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().Error(decodeResponse(resp))

	client := pb.NewAuthClient(conn)

	response, err := client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	marshalReq, err = json.Marshal(admin.SubnetRequest{IPWithMask: h.subnet.String()})
	h.Assert().NoError(err)
	url = fmt.Sprintf("http://%s/admin/remove_from_whitelist", h.cfg.HTTPServer.Addr)
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(marshalReq))
	h.Assert().NoError(err)

	defer resp.Body.Close()
	h.Assert().NoError(decodeResponse(resp))

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().True(response.Ok)

	response, err = client.SignIn(context.TODO(), &pb.SignInRequest{
		Login:    h.credentials.Login,
		Password: h.credentials.Password,
		Ip:       h.userIP.String(),
	})
	h.Assert().NoError(err)
	h.Assert().False(response.Ok)
}

func decodeResponse(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return errors.Wrap(shared.ErrEmpty, string(body))
	}

	return nil
}
