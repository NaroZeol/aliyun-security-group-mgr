package conf

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	if config == nil {
		t.Fatal("NewConfig returned nil")
	}
	if config.Credential == nil {
		t.Error("Credential is nil")
	}
	if config.Reloader == nil {
		t.Error("Reloader is nil")
	}
	if config.ECS == nil {
		t.Error("ECS is nil")
	}
	if config.SecurityGroup == nil {
		t.Error("SecurityGroup is nil")
	}
}

func TestLoadFile(t *testing.T) {
	// Create a temporary .env file
	tmpfile, err := os.CreateTemp("", "testenv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	content := []byte("ALIYUN_SGMGR_CREDENTIAL_TYPE=test_type")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading the file
	if err := LoadFile(tmpfile.Name()); err != nil {
		t.Errorf("LoadFile() error = %v", err)
	}

	// Verify env var is set
	if got := os.Getenv("ALIYUN_SGMGR_CREDENTIAL_TYPE"); got != "test_type" {
		t.Errorf("Env var not set correctly, got %v, want test_type", got)
	}
}

func TestLoadGlobalFromEnv(t *testing.T) {
	// Set some environment variables
	os.Setenv("ALIYUN_SGMGR_CREDENTIAL_TYPE", "env_type")
	os.Setenv("ALIYUN_SGMGR_CREDENTIAL_ACCESS_KEY_ID", "key_id")
	os.Setenv("ALIYUN_SGMGR_RELOADER_ENABLED", "true")
	os.Setenv("ALIYUN_SGMGR_RELOADER_INTERVAL", "60")

	defer func() {
		os.Unsetenv("ALIYUN_SGMGR_CREDENTIAL_TYPE")
		os.Unsetenv("ALIYUN_SGMGR_CREDENTIAL_ACCESS_KEY_ID")
		os.Unsetenv("ALIYUN_SGMGR_RELOADER_ENABLED")
		os.Unsetenv("ALIYUN_SGMGR_RELOADER_INTERVAL")
	}()

	config, err := LoadGlobalFromEnv()
	if err != nil {
		t.Fatalf("LoadGlobalFromEnv() error = %v", err)
	}

	if config.Credential.Type == nil || *config.Credential.Type != "env_type" {
		t.Errorf("Credential.Type = %v, want env_type", config.Credential.Type)
	}
	if config.Credential.AccessKeyId == nil || *config.Credential.AccessKeyId != "key_id" {
		t.Errorf("Credential.AccessKeyId = %v, want key_id", config.Credential.AccessKeyId)
	}
	if config.Reloader.Enabled == nil || *config.Reloader.Enabled != true {
		t.Errorf("Reloader.Enabled = %v, want true", config.Reloader.Enabled)
	}
	if config.Reloader.Interval == nil || *config.Reloader.Interval != 60 {
		t.Errorf("Reloader.Interval = %v, want 60", config.Reloader.Interval)
	}
}
