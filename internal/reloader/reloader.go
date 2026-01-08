package reloader

import (
	"aliyun-security-group-mgr/internal/conf"
	"log"
	"os"
	"time"
)

type Reloader struct {
	Config *conf.GlobalConfiguration

	reloadChan      chan struct{}
	expectedEntries []Entry
	lastReloadTime  time.Time
}

func NewReloader(config *conf.GlobalConfiguration, reloadChan chan struct{}) (*Reloader, error) {
	return &Reloader{
		Config:     config,
		reloadChan: reloadChan,
	}, nil
}

func (r *Reloader) Start() error {
	go r.watchEntries()
	select {}
}

func (r *Reloader) watchEntries() {
	ticker := time.NewTicker(time.Duration(*r.Config.Reloader.Interval) * time.Second)
	defer ticker.Stop()

	for {
		r.reloadEntries()
		<-ticker.C
	}
}

func (r *Reloader) reloadEntries() {
	// Check file modification time
	fileInfo, err := os.Stat(*r.Config.Reloader.WatchPath)
	if err != nil {
		log.Printf("[Reloader] failed to stat file: %v", err)
		return
	}
	if fileInfo.ModTime().Equal(r.lastReloadTime) {
		// No changes
		return
	}
	r.lastReloadTime = fileInfo.ModTime()

	// Read entries from file
	entries, err := ReadEntriesFromFile(*r.Config.Reloader.WatchPath)
	if err != nil {
		log.Printf("[Reloader] failed to read entries from file: %v", err)
		return
	}

	// Update expected entries
	r.expectedEntries = entries

	// Log reloading
	log.Printf("[Reloader] reloading rules from %s", *r.Config.Reloader.WatchPath)

	// Notify service to sync
	r.reloadChan <- struct{}{}
}

func (r *Reloader) GetExpectedEntries() []Entry {
	return r.expectedEntries
}
