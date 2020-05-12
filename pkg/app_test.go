package pkg

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/repository"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
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
	app         *App
}

func (g *GRPCBruteForceTestSuit) SetupSuite() {
	cfg, err := config.LoadFromFile("../init/config_test.yaml")
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

	g.app = NewApp(g.cfg,
		WithLeakyBuckets(g.cfg.RateLimiter),
		WithIPList(g.ipList),
	)

	go g.app.Run()
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
