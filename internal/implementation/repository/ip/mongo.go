package ip

import (
	"context"
	"errors"
	"fmt"
	"net"

	pkgerrors "github.com/pkg/errors"
	"github.com/teploff/antibruteforce/config"
	"github.com/teploff/antibruteforce/internal/domain/repository"
	"github.com/teploff/antibruteforce/internal/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoIPList struct {
	client    *mongo.Client
	whitelist *mongo.Collection
	blacklist *mongo.Collection
}

// NewMongoIPList returns mongodb repository of ip list.
func NewMongoIPList(cfg config.MongoConfig) (repository.IPStorable, error) {
	dns := fmt.Sprintf("mongodb://%s", cfg.Addr)
	client, err := mongo.NewClient(options.Client().ApplyURI(dns))

	if err != nil {
		return nil, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &mongoIPList{
		client:    client,
		whitelist: client.Database(cfg.DBName).Collection("whitelist"),
		blacklist: client.Database(cfg.DBName).Collection("blacklist"),
	}, nil
}

func (m *mongoIPList) AddInWhitelist(ipNet *net.IPNet) error {
	subnet := bson.M{"subnet": ipNet.String()}

	if err := m.whitelist.FindOne(context.TODO(), subnet).Err(); !errors.Is(err, mongo.ErrNoDocuments) && err != nil {
		return err
	} else if err == nil {
		return pkgerrors.Wrap(shared.ErrAlreadyExist, "in whitelist")
	}

	if err := m.blacklist.FindOne(context.TODO(), subnet).Err(); !errors.Is(err, mongo.ErrNoDocuments) && err != nil {
		return err
	} else if err == nil {
		return pkgerrors.Wrap(shared.ErrAlreadyExist, "in blacklist")
	}

	if _, err := m.whitelist.InsertOne(context.TODO(), subnet); err != nil {
		return err
	}

	return nil
}

func (m *mongoIPList) RemoveFromWhitelist(ipNet *net.IPNet) error {
	subnet := bson.M{"subnet": ipNet.String()}

	if err := m.whitelist.FindOneAndDelete(context.TODO(), subnet).Err(); err != nil {
		return err
	}

	return nil
}

func (m *mongoIPList) AddInBlacklist(ipNet *net.IPNet) error {
	subnet := bson.M{"subnet": ipNet.String()}

	if err := m.blacklist.FindOne(context.TODO(), subnet).Err(); !errors.Is(err, mongo.ErrNoDocuments) && err != nil {
		return err
	} else if err == nil {
		return pkgerrors.Wrap(shared.ErrAlreadyExist, "in blacklist")
	}

	if err := m.whitelist.FindOne(context.TODO(), subnet).Err(); !errors.Is(err, mongo.ErrNoDocuments) && err != nil {
		return err
	} else if err == nil {
		return pkgerrors.Wrap(shared.ErrAlreadyExist, "in whitelist")
	}

	if _, err := m.blacklist.InsertOne(context.TODO(), subnet); err != nil {
		return err
	}

	return nil
}

func (m *mongoIPList) RemoveFromBlacklist(ipNet *net.IPNet) error {
	subnet := bson.M{"subnet": ipNet.String()}

	if err := m.blacklist.FindOneAndDelete(context.TODO(), subnet).Err(); err != nil {
		return err
	}

	return nil
}

func (m *mongoIPList) IsIPInWhiteList(ip net.IP) (bool, error) {
	return m.isIPContains(ip, m.whitelist)
}

func (m *mongoIPList) IsIPInBlackList(ip net.IP) (bool, error) {
	return m.isIPContains(ip, m.blacklist)
}

func (m *mongoIPList) WhiteListLength() (int, error) {
	res, err := m.whitelist.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (m *mongoIPList) BlackListLength() (int, error) {
	res, err := m.blacklist.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (m *mongoIPList) Close() error {
	return m.client.Disconnect(context.TODO())
}

func (m *mongoIPList) isIPContains(ip net.IP, collection *mongo.Collection) (bool, error) {
	type Sub struct {
		Subnet string
	}

	cur, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		return false, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var sub Sub
		if err = cur.Decode(&sub); err != nil {
			return false, err
		}

		_, ipNet, err := net.ParseCIDR(sub.Subnet)
		if err != nil {
			return false, err
		}

		if ipNet.Contains(ip) {
			return true, nil
		}
	}

	if err = cur.Err(); err != nil {
		return false, err
	}

	return false, nil
}
