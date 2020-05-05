package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server      GrpcServer        `yaml:"server"`
	Logger      LoggerConfig      `yaml:"logger"`
	RateLimiter RateLimiterConfig `yaml:"rate_limiter"`
}

// GrpcServer configuration of grpc-instance service.
type GrpcServer struct {
	Addr string `yaml:"addr"`
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
	Login    Login    `yaml:"login"`
	Password Password `yaml:"password"`
	IP       IP       `yaml:"ip"`
	GCTime   int      `yaml:"gc_time"`
}

type Login struct {
	RPM        int `yaml:"rpm"`
	ExpireTime int `yaml:"expire_time"`
}

type Password struct {
	RPM        int `yaml:"rpm"`
	ExpireTime int `yaml:"expire_time"`
}

type IP struct {
	RPM        int `yaml:"rpm"`
	ExpireTime int `yaml:"expire_time"`
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
