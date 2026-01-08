package service

import (
	"aliyun-security-group-mgr/internal/conf"
	"aliyun-security-group-mgr/internal/ecs"
	"aliyun-security-group-mgr/internal/reloader"
)

type Service struct {
	Config   *conf.GlobalConfiguration
	Ecs      *ecs.Clerk
	Reloader *reloader.Reloader
}

func NewService(config *conf.GlobalConfiguration) (*Service, error) {
	return &Service{
		Config: config,
	}, nil
}

func (s *Service) Start() error {
	// New ECS Clerk
	ecsClerk, err := ecs.NewClerk(s.Config)
	if err != nil {
		return err
	}
	s.Ecs = ecsClerk

	// Check and create watch file if not exists
	err = s.checkWatchFile()
	if err != nil {
		return err
	}

	// New Reloader
	reloadChan := make(chan struct{})
	reloader, err := reloader.NewReloader(s.Config, reloadChan)
	if err != nil {
		return err
	}
	s.Reloader = reloader

	go s.Reloader.Start()

	for range reloadChan {
		s.syncSecurityGroupEntries()
	}
	return nil
}
