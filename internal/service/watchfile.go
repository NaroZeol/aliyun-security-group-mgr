package service

import (
	"aliyun-security-group-mgr/internal/reloader"

	"log"
	"os"
	"time"
)

func (s *Service) createNewWatchFile() error {
	f, err := os.Create(*s.Config.Reloader.WatchPath)
	if err != nil {
		log.Printf("[Service] failed to create watch file: %v", err)
		return err
	}
	f.Close()
	log.Printf("[Service] created watch file: %s", *s.Config.Reloader.WatchPath)
	log.Printf("[Service] fetching rules from ECS")
	currentEntries, err := s.getCurrentEntries()
	if err != nil {
		return err
	}
	for i := range currentEntries {
		currentEntries[i].ExpireAt, _ = time.Parse(time.RFC3339, "2100-01-01T00:00:00Z")
	}

	err = reloader.WriteEntriesToFile(*s.Config.Reloader.WatchPath, currentEntries)
	if err != nil {
		log.Printf("[Service] failed to write current rules to watch file: %v", err)
		return err
	}

	log.Printf("[Service] writing current rules to watch file")
	return nil
}

func (s *Service) checkWatchFile() error {
	_, err := os.Stat(*s.Config.Reloader.WatchPath)
	if err != nil {
		log.Printf("[Service] reloader watch path does not exist: %v", err)
		err = s.createNewWatchFile()
		if err != nil {
			return err
		}
	}
	return nil
}
