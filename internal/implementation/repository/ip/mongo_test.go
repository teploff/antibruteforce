package ip

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/internal/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoIPListTestSuit struct {
	dbName     string
	repository repository.IPStorable
	client     *mongo.Client
	suite.Suite
}

func TestMongoIPList(t *testing.T) {
	suite.Run(t, new(MongoIPListTestSuit))
}

func (m *MongoIPListTestSuit) SetupSuite() {
	cfg, _ := config.LoadFromFile("../../../../init/config_test.yaml")
	m.dbName = cfg.Mongo.DBName

	m.client, _ = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", cfg.Mongo.Addr)))
	_ = m.client.Connect(context.TODO())

	m.repository, _ = NewMongoIPList(cfg.Mongo)
}

func (m *MongoIPListTestSuit) TearDownSuite() {
	_ = m.repository.Close()
	_ = m.client.Disconnect(context.TODO())
}

func (m *MongoIPListTestSuit) TearDownTest() {
	_ = m.client.Database(m.dbName).Collection("whitelist").Drop(context.TODO())
	_ = m.client.Database(m.dbName).Collection("blacklist").Drop(context.TODO())
}

func (m *MongoIPListTestSuit) TestElementsAlreadyExistInWhiteAndBlacklists() {
	_, whiteNet, _ := net.ParseCIDR("192.168.128.0/24")
	_, blackNet, _ := net.ParseCIDR("192.168.131.0/24")

	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet))
	m.Assert().Error(m.repository.AddInWhitelist(whiteNet))

	m.Assert().NoError(m.repository.AddInBlacklist(blackNet))
	m.Assert().Error(m.repository.AddInBlacklist(blackNet))

	length, err := m.repository.WhiteListLength()
	m.Assert().Equal(1, length)
	m.Assert().NoError(err)

	length, err = m.repository.BlackListLength()
	m.Assert().Equal(1, length)
	m.Assert().NoError(err)
}

func (m *MongoIPListTestSuit) TestElementInWhiteAndBlacklistsSimultaneously() {
	_, ipNet, _ := net.ParseCIDR("192.168.128.0/24")

	m.Assert().NoError(m.repository.AddInWhitelist(ipNet))
	m.Assert().Error(m.repository.AddInBlacklist(ipNet))

	length, err := m.repository.WhiteListLength()
	m.Assert().Equal(1, length)
	m.Assert().NoError(err)

	length, err = m.repository.BlackListLength()
	m.Assert().Equal(0, length)
	m.Assert().NoError(err)

	m.Assert().NoError(m.repository.RemoveFromWhitelist(ipNet))

	m.Assert().NoError(m.repository.AddInBlacklist(ipNet))
	m.Assert().Error(m.repository.AddInWhitelist(ipNet))

	length, err = m.repository.WhiteListLength()
	m.Assert().Equal(0, length)
	m.Assert().NoError(err)

	length, err = m.repository.BlackListLength()
	m.Assert().Equal(1, length)
	m.Assert().NoError(err)
}

func (m *MongoIPListTestSuit) TestRemovingWhiteAndBlackLists() {
	_, whiteNet1, _ := net.ParseCIDR("192.168.128.0/24")
	_, whiteNet2, _ := net.ParseCIDR("192.168.129.0/24")
	_, whiteNet3, _ := net.ParseCIDR("192.168.130.0/24")

	_, blackNet1, _ := net.ParseCIDR("192.168.131.0/24")
	_, blackNet2, _ := net.ParseCIDR("192.168.132.0/24")
	_, blackNet3, _ := net.ParseCIDR("192.168.133.0/24")

	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet1))
	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet2))
	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet3))

	m.Assert().NoError(m.repository.AddInBlacklist(blackNet1))
	m.Assert().NoError(m.repository.AddInBlacklist(blackNet2))
	m.Assert().NoError(m.repository.AddInBlacklist(blackNet3))

	m.Assert().NoError(m.repository.RemoveFromWhitelist(whiteNet1))
	m.Assert().NoError(m.repository.RemoveFromWhitelist(whiteNet2))
	m.Assert().NoError(m.repository.RemoveFromWhitelist(whiteNet3))

	m.Assert().NoError(m.repository.RemoveFromBlacklist(blackNet1))
	m.Assert().NoError(m.repository.RemoveFromBlacklist(blackNet2))
	m.Assert().NoError(m.repository.RemoveFromBlacklist(blackNet3))

	length, err := m.repository.WhiteListLength()
	m.Assert().Equal(0, length)
	m.Assert().NoError(err)
	length, err = m.repository.BlackListLength()
	m.Assert().Equal(0, length)
	m.Assert().NoError(err)
}

func (m *MongoIPListTestSuit) TestBelongWhiteAndBlackLists() {
	_, whiteNet1, _ := net.ParseCIDR("192.168.130.0/24")
	_, whiteNet2, _ := net.ParseCIDR("192.168.0.0/16")
	_, whiteNet3, _ := net.ParseCIDR("192.0.0.0/8")

	_, blackNet1, _ := net.ParseCIDR("10.200.128.0/24")
	_, blackNet2, _ := net.ParseCIDR("10.200.0.0/16")
	_, blackNet3, _ := net.ParseCIDR("10.0.0.0/8")

	whiteIP1 := net.ParseIP("192.15.10.11")
	whiteIP2 := net.ParseIP("192.168.10.11")
	whiteIP3 := net.ParseIP("192.168.130.11")

	blackIP1 := net.ParseIP("10.15.10.11")
	blackIP2 := net.ParseIP("10.200.10.11")
	blackIP3 := net.ParseIP("10.200.128.11")

	neutralIP := net.ParseIP("127.0.0.1")

	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet1))
	exist, err := m.repository.IsIPInWhiteList(whiteIP1)
	m.Assert().False(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInWhiteList(whiteIP2)
	m.Assert().False(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInWhiteList(whiteIP3)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	m.Assert().NoError(m.repository.RemoveFromWhitelist(whiteNet1))

	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet2))
	exist, err = m.repository.IsIPInWhiteList(whiteIP1)
	m.Assert().False(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInWhiteList(whiteIP2)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInWhiteList(whiteIP3)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	m.Assert().NoError(m.repository.RemoveFromWhitelist(whiteNet2))

	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet3))

	exist, err = m.repository.IsIPInWhiteList(whiteIP1)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInWhiteList(whiteIP2)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInWhiteList(whiteIP3)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	m.Assert().NoError(m.repository.RemoveFromWhitelist(whiteNet3))

	m.Assert().NoError(m.repository.AddInBlacklist(blackNet1))
	exist, err = m.repository.IsIPInBlackList(blackIP1)
	m.Assert().False(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInBlackList(blackIP2)
	m.Assert().False(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInBlackList(blackIP3)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	m.Assert().NoError(m.repository.RemoveFromBlacklist(blackNet1))

	m.Assert().NoError(m.repository.AddInBlacklist(blackNet2))
	exist, err = m.repository.IsIPInBlackList(blackIP1)
	m.Assert().False(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInBlackList(blackIP2)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInBlackList(blackIP3)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	m.Assert().NoError(m.repository.RemoveFromBlacklist(blackNet2))

	m.Assert().NoError(m.repository.AddInBlacklist(blackNet3))
	exist, err = m.repository.IsIPInBlackList(blackIP1)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInBlackList(blackIP2)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInBlackList(blackIP3)
	m.Assert().True(exist)
	m.Assert().NoError(err)
	m.Assert().NoError(m.repository.RemoveFromBlacklist(blackNet3))

	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet1))
	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet2))
	m.Assert().NoError(m.repository.AddInWhitelist(whiteNet3))

	m.Assert().NoError(m.repository.AddInBlacklist(blackNet1))
	m.Assert().NoError(m.repository.AddInBlacklist(blackNet2))
	m.Assert().NoError(m.repository.AddInBlacklist(blackNet3))

	exist, err = m.repository.IsIPInWhiteList(neutralIP)
	m.Assert().False(exist)
	m.Assert().NoError(err)
	exist, err = m.repository.IsIPInBlackList(neutralIP)
	m.Assert().False(exist)
	m.Assert().NoError(err)

	length, err := m.repository.WhiteListLength()
	m.Assert().Equal(3, length)
	m.Assert().NoError(err)
	length, err = m.repository.BlackListLength()
	m.Assert().Equal(3, length)
	m.Assert().NoError(err)
}
