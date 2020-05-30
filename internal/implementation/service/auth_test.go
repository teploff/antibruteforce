package service

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/internal/domain/entity"
	"github.com/teploff/antibruteforce/internal/domain/repository"
	"github.com/teploff/antibruteforce/internal/implementation/repository/bucket"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
	"github.com/teploff/antibruteforce/internal/limiter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceTestSuit struct {
	suite.Suite
	cfg         config.Config
	client      *mongo.Client
	credentials entity.Credentials
	userIP      net.IP
	subNet      *net.IPNet
	login       repository.BucketStorable
	password    repository.BucketStorable
	ip          repository.BucketStorable
	ipList      repository.IPStorable
	rateLimiter *limiter.RateLimiter
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuit))
}

func (a *AuthServiceTestSuit) SetupSuite() {
	cfg, _ := config.LoadFromFile("../../../init/config_test.yaml")
	a.cfg = cfg
	a.cfg.Mongo.DBName = "test_auth"

	a.login = bucket.NewLeakyBucket(a.cfg.RateLimiter.Login.Rate, a.cfg.RateLimiter.Login.Interval,
		a.cfg.RateLimiter.Login.ExpireTime)
	a.password = bucket.NewLeakyBucket(a.cfg.RateLimiter.Password.Rate, a.cfg.RateLimiter.Password.Interval,
		a.cfg.RateLimiter.Password.ExpireTime)
	a.ip = bucket.NewLeakyBucket(a.cfg.RateLimiter.IP.Rate, a.cfg.RateLimiter.IP.Interval, a.cfg.RateLimiter.IP.ExpireTime)
	a.ipList, _ = ip.NewMongoIPList(a.cfg.Mongo)
	a.rateLimiter = limiter.NewRateLimiter(a.login, a.password, a.ip, a.cfg.RateLimiter.GCTime)

	a.client, _ = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", a.cfg.Mongo.Addr)))
	_ = a.client.Connect(context.TODO())

	a.credentials = entity.Credentials{
		Login:    "loginAuth",
		Password: "passwordAuth",
	}
	a.userIP = net.ParseIP("192.168.199.132")
	_, a.subNet, _ = net.ParseCIDR("192.168.199.0/24")
}

func (a *AuthServiceTestSuit) TearDownSuite() {
	_ = a.ipList.Close()
	_ = a.client.Disconnect(context.TODO())
}

func (a *AuthServiceTestSuit) TearDownTest() {
	_ = a.client.Database(a.cfg.Mongo.DBName).Collection("whitelist").Drop(context.TODO())
	_ = a.client.Database(a.cfg.Mongo.DBName).Collection("blacklist").Drop(context.TODO())
	a.login = bucket.NewLeakyBucket(a.cfg.RateLimiter.Login.Rate, a.cfg.RateLimiter.Login.Interval,
		a.cfg.RateLimiter.Login.ExpireTime)
	a.password = bucket.NewLeakyBucket(a.cfg.RateLimiter.Password.Rate, a.cfg.RateLimiter.Password.Interval,
		a.cfg.RateLimiter.Password.ExpireTime)
	a.ip = bucket.NewLeakyBucket(a.cfg.RateLimiter.IP.Rate, a.cfg.RateLimiter.IP.Interval,
		a.cfg.RateLimiter.IP.ExpireTime)
	a.rateLimiter = limiter.NewRateLimiter(a.login, a.password, a.ip, a.cfg.RateLimiter.GCTime)
}

func (a *AuthServiceTestSuit) TestIPInWhitelist() {
	authSvc := NewAuthService(a.rateLimiter, a.ipList)
	a.Assert().NoError(a.ipList.AddInWhitelist(a.subNet))

	allow, err := authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
}

func (a *AuthServiceTestSuit) TestIPInBlacklist() {
	authSvc := NewAuthService(a.rateLimiter, a.ipList)
	a.Assert().NoError(a.ipList.AddInBlacklist(a.subNet))

	allow, err := authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().False(allow)

	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().False(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().False(allow)
}

func (a *AuthServiceTestSuit) TestBruteForceByLogin() {
	authSvc := NewAuthService(a.rateLimiter, a.ipList)
	allow, err := authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().False(allow)
}

func (a *AuthServiceTestSuit) TestBruteForceByPassword() {
	authSvc := NewAuthService(a.rateLimiter, a.ipList)
	allow, err := authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().False(allow)
}

func (a *AuthServiceTestSuit) TestBruteForceByIP() {
	authSvc := NewAuthService(a.rateLimiter, a.ipList)
	allow, err := authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().True(allow)
	allow, err = authSvc.LogIn(a.credentials, a.userIP)
	a.Assert().NoError(err)
	a.Assert().False(allow)
}
