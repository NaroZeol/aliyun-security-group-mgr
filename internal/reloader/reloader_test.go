package reloader

import (
	"aliyun-security-group-mgr/internal/conf"
	"os"
	"testing"
	"time"
)

func TestReloader_reloadEntries(t *testing.T) {
	// Create temp file
	tmpfile, err := os.CreateTemp("", "rules.conf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write initial content
	initialContent := "accept ingress tcp 22/22 from 0.0.0.0/0 priority 100 until 2024-12-31T23:59:59+08:00 # SSH\n"
	if err := os.WriteFile(tmpfile.Name(), []byte(initialContent), 0644); err != nil {
		t.Fatal(err)
	}

	watchPath := tmpfile.Name()
	interval := int64(10)
	config := &conf.GlobalConfiguration{
		Reloader: &conf.Reloader{
			WatchPath: &watchPath,
			Interval:  &interval,
		},
	}

	reloadChan := make(chan struct{}, 1)
	reloader, err := NewReloader(config, reloadChan)
	if err != nil {
		t.Fatalf("NewReloader error = %v", err)
	}

	// 1. First reload
	reloader.reloadEntries()

	// Check channel
	select {
	case <-reloadChan:
		// success
	default:
		t.Error("expected signal in reloadChan")
	}

	entries := reloader.GetExpectedEntries()
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].SecurityGroup.Description != "SSH" {
		t.Errorf("expected description SSH, got %s", entries[0].SecurityGroup.Description)
	}

	// 2. Second reload with NO change
	reloader.reloadEntries()
	select {
	case <-reloadChan:
		t.Error("unexpected signal in reloadChan (file not changed)")
	default:
		// success
	}

	// 3. Update file and reload
	time.Sleep(1 * time.Second)

	newContent := "accept ingress tcp 80/80 from 0.0.0.0/0 priority 100 until 2024-12-31T23:59:59+08:00 # HTTP\n"
	if err := os.WriteFile(tmpfile.Name(), []byte(newContent), 0644); err != nil {
		t.Fatal(err)
	}

	reloader.reloadEntries()

	select {
	case <-reloadChan:
		// success
	default:
		t.Error("expected signal in reloadChan after file update")
	}

	entries = reloader.GetExpectedEntries()
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].SecurityGroup.Description != "HTTP" {
		t.Errorf("expected description HTTP, got %s", entries[0].SecurityGroup.Description)
	}
}
