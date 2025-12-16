package reloader

import (
	"aliyun-security-group-mgr/internal/ecs"
	"aliyun-security-group-mgr/internal/utils"

	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

type Entry struct {
	ecs.SecurityGroupRule
	ExpireAt time.Time
}

func ReadEntriesFromFile(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var entries []Entry
	for _, line := range lines {
		entry, err := DecodeEntry(line)
		if err != nil {
			return nil, err
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func DecodeEntry(line string) (*Entry, error) {
	comment := utils.ExtractCommentFromLine(line)
	line = utils.RemoveCommentFromLine(line)
	parts := strings.Fields(line)
	if len(parts) != 10 {
		return nil, fmt.Errorf("invalid entry line: %s", line)
	}

	policy := ecs.Policy(strings.Title(parts[0]))
	direction := parts[1]
	ipProtocol := strings.ToUpper(parts[2])
	portRange := parts[3]
	cidrIp := parts[5]
	priority := parts[7]
	expireAtStr := parts[9]

	expireAt, err := time.Parse(time.RFC3339, expireAtStr)
	if err != nil {
		return nil, fmt.Errorf("invalid expire at format: %s", expireAtStr)
	}

	entry := &Entry{
		SecurityGroupRule: ecs.SecurityGroupRule{
			Policy:      policy,
			Direction:   direction,
			IpProtocol:  ipProtocol,
			PortRange:   portRange,
			CidrIp:      cidrIp,
			Priority:    priority,
			Description: comment,
		},
		ExpireAt: expireAt,
	}

	return entry, nil
}

func WriteEntriesToFile(path string, entries []Entry) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range entries {
		line := EncodeEntry(entry)
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func EncodeEntry(entry Entry) string {
	var policy string
	{
		runes := []rune(string(entry.Policy))
		runes[0] = unicode.ToLower(runes[0])
		policy = string(runes)
	}
	var direction string = entry.Direction
	var ipProtocol string = strings.ToLower(entry.IpProtocol)
	var portRange string = entry.PortRange
	var directionWord string
	if entry.Direction == "ingress" {
		directionWord = "from"
	} else {
		directionWord = "to"
	}
	var cidrIp string = entry.CidrIp
	var priority string = entry.Priority
	var expireAt string = entry.ExpireAt.Format(time.RFC3339)

	str := fmt.Sprintf("%s %s %s %s %s %s priority %s until %s",
		policy,
		direction,
		ipProtocol,
		portRange,
		directionWord,
		cidrIp,
		priority,
		expireAt,
	)

	str = strings.TrimSpace(str)
	if entry.Description != "" {
		str += " # " + entry.Description
	}

	return str
}
