package service

import (
	"aliyun-security-group-mgr/internal/conf"
	"flag"
)

var (
	configFile string
)

func Init() (*Service, error) {
	flag.StringVar(&configFile, "config", ".env", "Path to configuration file")
	flag.Parse()

	if err := conf.LoadFile(configFile); err != nil {
		return nil, err
	}

	config, err := conf.LoadGlobalFromEnv()
	if err != nil {
		return nil, err
	}

	service, err := NewService(config)
	if err != nil {
		return nil, err
	}

	return service, err
}
