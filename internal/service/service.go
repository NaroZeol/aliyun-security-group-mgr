package service

import (
	"aliyun-security-group-mgr/internal/conf"
	"aliyun-security-group-mgr/internal/ecs"
)

type Service struct {
	Config *conf.GlobalConfiguration
	Ecs    *ecs.Clerk
}

func NewService(config *conf.GlobalConfiguration) (*Service, error) {
	ecsClerk, err := ecs.NewClerk(config)
	if err != nil {
		return nil, err
	}

	return &Service{
		Config: config,
		Ecs:    ecsClerk,
	}, nil
}

func (s *Service) Start() error {
	select {}
}
