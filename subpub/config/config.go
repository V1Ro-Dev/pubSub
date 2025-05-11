package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	addr                 = ":8080"
	maxRecvMsgSize       = 10485760
	maxSendMsgSize       = 10485760
	maxConcurrentStreams = 100
)

var defaultCfg = &Config{
	Addr:                 addr,
	MaxRecvMsgSize:       maxRecvMsgSize,
	MaxSendMsgSize:       maxSendMsgSize,
	MaxConcurrentStreams: maxConcurrentStreams,
}

type Config struct {
	Addr                 string `yaml:"addr"`
	MaxRecvMsgSize       int    `yaml:"max_recv_msg_size"`
	MaxSendMsgSize       int    `yaml:"max_send_msg_size"`
	MaxConcurrentStreams uint32 `yaml:"max_concurrent_streams"`
}

func loadConfig(path string) (*Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return defaultCfg, err
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return defaultCfg, err
	}

	return &cfg, nil
}

func Parse(configPath string) (*Config, error) {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return cfg, fmt.Errorf("load config error: %w", err)
	}

	return cfg, nil
}
