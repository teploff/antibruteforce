package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Config holds all configs.
type Config struct {
	GRPCServer  GRPCConfig        `yaml:"gRPC_server"`
	HTTPServer  HTTPConfig        `yaml:"http_server"`
	Mongo       MongoConfig       `yaml:"mongo"`
	Logger      LoggerConfig      `yaml:"logger"`
	RateLimiter RateLimiterConfig `yaml:"rate_limiter"`
}

// GRPCConfig configuration of gRPC-instance service.
type GRPCConfig struct {
	Addr string `yaml:"addr"`
}

// HTTPConfig configuration of http-instance service.
type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

// MongoConfig configuration of mongoDB database.
type MongoConfig struct {
	Addr   string `yaml:"addr"`
	DBName string `yaml:"db_name"`
}

// LoggerConfig logger configuration.
//
// Filename - log file name.
//
// MaxSize - max log file size.
type LoggerConfig struct {
	Filename string `yaml:"file_name"`
	MaxSize  int    `yaml:"max_size"`
	Level    string `yaml:"level"`
}

// RateLimiterConfig holds all leaky bucket configs
//
// GCTime - launch time garbage collector which delete expired buckets.
type RateLimiterConfig struct {
	Login    Login         `yaml:"login"`
	Password Password      `yaml:"password"`
	IP       IP            `yaml:"ip"`
	GCTime   time.Duration `yaml:"gc_time"`
}

// Login leaky bucket config for logins. Rate/Interval requests per login bucket with expired time ExpireTime.
type Login struct {
	Rate       int           `yaml:"rate"`
	Interval   time.Duration `yaml:"interval"`
	ExpireTime time.Duration `yaml:"expire_time"`
}

// Password leaky bucket config for password. Rate/Interval requests per password bucket with expired time ExpireTime.
type Password struct {
	Rate       int           `yaml:"rate"`
	Interval   time.Duration `yaml:"interval"`
	ExpireTime time.Duration `yaml:"expire_time"`
}

// IP leaky bucket config for ip. Rate/Interval requests per ip bucket with expired time ExpireTime.
type IP struct {
	Rate       int           `yaml:"rate"`
	Interval   time.Duration `yaml:"interval"`
	ExpireTime time.Duration `yaml:"expire_time"`
}

// LoadFromFile create configuration from file.
func LoadFromFile(fileName string) (Config, error) {
	cfg := Config{}

	configBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
