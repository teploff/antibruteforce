package service

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/domain/entity"
	"github.com/teploff/antibruteforce/domain/repository"
	"github.com/teploff/antibruteforce/internal/implementation/repository/bucket"
	"github.com/teploff/antibruteforce/internal/implementation/repository/ip"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminServiceTestSuit struct {
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
}

func TestAdminService(t *testing.T) {
	suite.Run(t, new(AdminServiceTestSuit))
}

func (a *AdminServiceTestSuit) SetupSuite() {
	cfg, _ := config.LoadFromFile("../../../init/config_test.yaml")
	a.cfg = cfg
	a.cfg.Mongo.DBName = "test_admin"

	a.login = bucket.NewLeakyBucket(a.cfg.RateLimiter.Login.Rate, a.cfg.RateLimiter.Login.Interval,
		a.cfg.RateLimiter.Login.ExpireTime)
	a.password = bucket.NewLeakyBucket(a.cfg.RateLimiter.Password.Rate, a.cfg.RateLimiter.Password.Interval,
		a.cfg.RateLimiter.Password.ExpireTime)
	a.ip = bucket.NewLeakyBucket(a.cfg.RateLimiter.IP.Rate, a.cfg.RateLimiter.IP.Interval, a.cfg.RateLimiter.IP.ExpireTime)
	a.ipList, _ = ip.NewMongoIPList(a.cfg.Mongo)

	a.client, _ = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", a.cfg.Mongo.Addr)))
	_ = a.client.Connect(context.TODO())

	a.credentials = entity.Credentials{
		Login:    "loginAuth",
		Password: "passwordAuth",
	}
	a.userIP = net.ParseIP("192.168.199.132")
	_, a.subNet, _ = net.ParseCIDR("192.168.199.0/24")
}

func (a *AdminServiceTestSuit) TearDownSuite() {
	_ = a.ipList.Close()
	_ = a.client.Disconnect(context.TODO())
}

func (a *AdminServiceTestSuit) TearDownTest() {
	_ = a.client.Database(a.cfg.Mongo.DBName).Collection("whitelist").Drop(context.TODO())
	_ = a.client.Database(a.cfg.Mongo.DBName).Collection("blacklist").Drop(context.TODO())
	a.login = bucket.NewLeakyBucket(a.cfg.RateLimiter.Login.Rate, a.cfg.RateLimiter.Login.Interval,
		a.cfg.RateLimiter.Login.ExpireTime)
	a.password = bucket.NewLeakyBucket(a.cfg.RateLimiter.Password.Rate, a.cfg.RateLimiter.Password.Interval,
		a.cfg.RateLimiter.Password.ExpireTime)
	a.ip = bucket.NewLeakyBucket(a.cfg.RateLimiter.IP.Rate, a.cfg.RateLimiter.IP.Interval,
		a.cfg.RateLimiter.IP.ExpireTime)
}

func (a *AdminServiceTestSuit) TestResetBucketByLogin() {
	adminSvc := NewAdminService(a.ipList, a.login, a.password, a.ip)
	_, err := a.login.Add(a.credentials.Login)
	a.Assert().NoError(err)

	a.Assert().NoError(adminSvc.ResetBucketByLogin(a.credentials.Login))
	a.Assert().Error(adminSvc.ResetBucketByLogin(a.credentials.Login))
}

func (a *AdminServiceTestSuit) TestResetBucketByPassword() {
	adminSvc := NewAdminService(a.ipList, a.login, a.password, a.ip)
	_, err := a.password.Add(a.credentials.Password)
	a.Assert().NoError(err)

	a.Assert().NoError(adminSvc.ResetBucketByPassword(a.credentials.Password))
	a.Assert().Error(adminSvc.ResetBucketByPassword(a.credentials.Password))
}

func (a *AdminServiceTestSuit) TestResetBucketByIP() {
	adminSvc := NewAdminService(a.ipList, a.login, a.password, a.ip)
	_, err := a.ip.Add(a.userIP.String())
	a.Assert().NoError(err)
	a.Assert().NoError(adminSvc.ResetBucketByIP(a.userIP))
	a.Assert().Error(adminSvc.ResetBucketByIP(a.userIP))
}

func (a *AdminServiceTestSuit) TestAddInBlacklist() {
	adminSvc := NewAdminService(a.ipList, a.login, a.password, a.ip)
	a.Assert().NoError(adminSvc.AddInBlacklist(a.subNet))
	a.Assert().Error(adminSvc.AddInBlacklist(a.subNet))
}

func (a *AdminServiceTestSuit) TestRemoveFromBlacklist() {
	adminSvc := NewAdminService(a.ipList, a.login, a.password, a.ip)
	a.Assert().NoError(adminSvc.AddInBlacklist(a.subNet))
	a.Assert().NoError(adminSvc.RemoveFromBlacklist(a.subNet))
	a.Assert().Error(adminSvc.RemoveFromBlacklist(a.subNet))
}

func (a *AdminServiceTestSuit) TestAddInWhitelist() {
	adminSvc := NewAdminService(a.ipList, a.login, a.password, a.ip)
	a.Assert().NoError(adminSvc.AddInWhitelist(a.subNet))
	a.Assert().Error(adminSvc.AddInWhitelist(a.subNet))
}

func (a *AdminServiceTestSuit) TestRemoveFromWhitelist() {
	adminSvc := NewAdminService(a.ipList, a.login, a.password, a.ip)
	a.Assert().NoError(adminSvc.AddInWhitelist(a.subNet))
	a.Assert().NoError(adminSvc.RemoveFromWhitelist(a.subNet))
	a.Assert().Error(adminSvc.RemoveFromWhitelist(a.subNet))
}
