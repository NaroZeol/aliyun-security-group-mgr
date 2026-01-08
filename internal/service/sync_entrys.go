package service

import (
	"aliyun-security-group-mgr/internal/reloader"

	"log"
	"time"
)

func (s *Service) getCurrentEntries() ([]reloader.Entry, error) {
	securityRule, err := s.Ecs.DescribeSecurityGroupAttribute()
	if err != nil {
		log.Printf("[Service] failed to get current entries: %v", err)
		return nil, err
	}
	var entries []reloader.Entry
	for _, rule := range securityRule {
		entry := reloader.Entry{
			SecurityGroup: rule,
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func buildMap(entries []reloader.Entry) map[string]reloader.Entry {
	buildKey := func(entry reloader.Entry) string {
		return entry.SecurityGroup.CidrIp + "|" + entry.SecurityGroup.IpProtocol + "|" + entry.SecurityGroup.PortRange + "|" + entry.SecurityGroup.Direction
	}

	result := make(map[string]reloader.Entry)
	for _, entry := range entries {
		key := buildKey(entry)
		result[key] = entry
	}
	return result
}

func (s *Service) syncSecurityGroupEntries() error {
	expectedEntries := s.Reloader.GetExpectedEntries()
	currentEntries, err := s.getCurrentEntries()
	if err != nil {
		return err
	}

	expectedEntriesMap := buildMap(expectedEntries)
	currentEntriesMap := buildMap(currentEntries)

	now := time.Now()
	entriesToAdd := []string{}
	entriesToUpdate := []string{}
	entriesToDelete := []string{}

	// Determine entries to add, update, delete
	for key, expectedEntry := range expectedEntriesMap {
		isExpired := expectedEntry.ExpireAt.Before(now)
		currentEntry, exists := currentEntriesMap[key]

		// no existing and not expired -> add
		if !exists && !isExpired {
			entriesToAdd = append(entriesToAdd, key)
			continue
		}

		// existing and not expired but different content -> modify
		if exists && !isExpired && !expectedEntry.EqualContent(currentEntry) {
			entriesToUpdate = append(entriesToUpdate, key)
			continue
		}

		// existing and expired -> delete
		if exists && isExpired {
			entriesToDelete = append(entriesToDelete, key)
			continue
		}
	}

	for key := range currentEntriesMap {
		_, exists := expectedEntriesMap[key]

		// existing in current but not in expected -> delete
		if !exists {
			entriesToDelete = append(entriesToDelete, key)
		}
	}

	log.Printf("[Service] synchronizing - to add: %d, to update: %d, to delete: %d", len(entriesToAdd), len(entriesToUpdate), len(entriesToDelete))

	// Additions
	for _, key := range entriesToAdd {
		entry := expectedEntriesMap[key]
		err := s.Ecs.AddSecurityGroupRule(entry.SecurityGroup)
		if err != nil {
			log.Printf("[Service] failed to add rule: %+v, error: %v", entry.SecurityGroup, err)
		} else {
			log.Printf("[Service] successfully added rule: %+v", entry.SecurityGroup)
		}
	}

	// Updates
	for _, key := range entriesToUpdate {
		oldEntry := currentEntriesMap[key]
		newEntry := expectedEntriesMap[key]
		err := s.Ecs.ModifySecurityGroupRule(oldEntry.SecurityGroup.Id, newEntry.SecurityGroup)
		if err != nil {
			log.Printf("[Service] failed to update rule from: %+v to: %+v, error: %v", oldEntry.SecurityGroup, newEntry.SecurityGroup, err)
		} else {
			log.Printf("[Service] successfully updated rule from: %+v to: %+v", oldEntry.SecurityGroup, newEntry.SecurityGroup)
		}
	}

	// Deletions
	for _, key := range entriesToDelete {
		entry := currentEntriesMap[key]
		err := s.Ecs.RemoveSecurityGroupRule(entry.SecurityGroup)
		if err != nil {
			log.Printf("[Service] failed to delete rule: %+v, error: %v", entry.SecurityGroup, err)
		} else {
			log.Printf("[Service] successfully deleted rule: %+v", entry.SecurityGroup)
		}
	}

	log.Printf("[Service] synchronization completed")

	return nil
}
