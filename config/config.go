package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GRPCServer  GRPCConfig        `yaml:"gRPC_server"`
	HTTPServer  HTTPConfig        `yaml:"http_server"`
	Redis       RedisConfig       `yaml:"redis"`
	Logger      LoggerConfig      `yaml:"logger"`
	RateLimiter RateLimiterConfig `yaml:"rate_limiter"`
}

// GRPCConfig configuration of grpc-instance service.
type GRPCConfig struct {
	Addr string `yaml:"addr"`
}

// HTTPConfig configuration of http-instance service.
type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

// RedisConfig configuration of cache database.
type RedisConfig struct {
	Addr        string `yaml:"addr"`
	Password    string `yaml:"password"`
	DbWhitelist int    `yaml:"db_whitelist"`
	DbBlacklist int    `yaml:"db_blacklist"`
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

type RateLimiterConfig struct {
	Login    Login         `yaml:"login"`
	Password Password      `yaml:"password"`
	IP       IP            `yaml:"ip"`
	GCTime   time.Duration `yaml:"gc_time"`
}

type Login struct {
	Rate       int           `yaml:"rate"`
	Interval   time.Duration `yaml:"interval"`
	ExpireTime time.Duration `yaml:"expire_time"`
}

type Password struct {
	Rate       int           `yaml:"rate"`
	Interval   time.Duration `yaml:"interval"`
	ExpireTime time.Duration `yaml:"expire_time"`
}

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
