package main

import (
	"aliyun-security-group-mgr/internal/conf"
	"aliyun-security-group-mgr/internal/service"
	"flag"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.Parse()

	if err := conf.LoadFile(configFile); err != nil {
		panic(err)
	}

	config, err := conf.LoadGlobalFromEnv()
	if err != nil {
		panic(err)
	}

	service, err := service.NewService(config)
	if err != nil {
		panic(err)
	}

	if err := service.Start(); err != nil {
		panic(err)
	}
}
