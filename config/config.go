package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server  GrpcServer   `yaml:"server"`
	Logger  LoggerConfig `yaml:"logger"`
	Limiter RateLimiter  `yaml:"limiter"`
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

type RateLimiter struct {
	Login    int `yaml:"login"`
	Password int `yaml:"password"`
	IP       int `yaml:"ip"`
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
