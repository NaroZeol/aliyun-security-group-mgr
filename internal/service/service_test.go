package service

import (
	"aliyun-security-group-mgr/internal/conf"
	"aliyun-security-group-mgr/internal/ecs"
	"aliyun-security-group-mgr/internal/reloader"
	"testing"
)

func TestNewService(t *testing.T) {
	config := &conf.GlobalConfiguration{}
	service, err := NewService(config)
	if err != nil {
		t.Fatalf("NewService error = %v", err)
	}
	if service.Config != config {
		t.Error("Config not set")
	}
}

func Test_buildMap(t *testing.T) {
	e1 := reloader.Entry{
		SecurityGroup: ecs.SecurityGroupRule{
			CidrIp:     "0.0.0.0/0",
			IpProtocol: "tcp",
			PortRange:  "80/80",
			Direction:  "ingress",
		},
	}
	e2 := reloader.Entry{
		SecurityGroup: ecs.SecurityGroupRule{
			CidrIp:     "1.2.3.4/32", // Diff Cidr
			IpProtocol: "tcp",
			PortRange:  "80/80",
			Direction:  "ingress",
		},
	}
	e3 := reloader.Entry{
		SecurityGroup: ecs.SecurityGroupRule{ // Same as e1
			CidrIp:     "0.0.0.0/0",
			IpProtocol: "tcp",
			PortRange:  "80/80",
			Direction:  "ingress",
		},
	}

	entries := []reloader.Entry{e1, e2}
	m := buildMap(entries)

	if len(m) != 2 {
		t.Errorf("buildMap count = %d, want 2", len(m))
	}

	// key check
	key1 := "0.0.0.0/0|tcp|80/80|ingress"
	if _, ok := m[key1]; !ok {
		t.Errorf("key %s not found", key1)
	}

	// overwriting check
	entriesOverwrite := []reloader.Entry{e1, e3}
	mOver := buildMap(entriesOverwrite)
	if len(mOver) != 1 {
		t.Errorf("buildMap count = %d, want 1 (overwrite)", len(mOver))
	}
}
