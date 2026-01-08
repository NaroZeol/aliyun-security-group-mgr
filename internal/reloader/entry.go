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
	SecurityGroup ecs.SecurityGroupRule
	ExpireAt      time.Time
}

func (e *Entry) EqualContent(other Entry) bool {
	return true &&
		e.SecurityGroup.CidrIp == other.SecurityGroup.CidrIp &&
		e.SecurityGroup.PortRange == other.SecurityGroup.PortRange &&
		e.SecurityGroup.IpProtocol == other.SecurityGroup.IpProtocol &&
		e.SecurityGroup.Policy == other.SecurityGroup.Policy &&
		e.SecurityGroup.Priority == other.SecurityGroup.Priority &&
		e.SecurityGroup.Direction == other.SecurityGroup.Direction &&
		e.SecurityGroup.Description == other.SecurityGroup.Description
}

func ReadEntriesFromFile(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
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
			if strings.Contains(err.Error(), "empty line") {
				continue
			}
			return nil, err
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func DecodeEntry(line string) (*Entry, error) {
	comment := utils.ExtractCommentFromLine(line)
	line = utils.RemoveCommentFromLine(line)
	if strings.TrimSpace(line) == "" {
		return nil, fmt.Errorf("empty line")
	}

	parts := strings.Fields(line)
	if len(parts) != 10 {
		return nil, fmt.Errorf("invalid entry line: %s", line)
	}

	policy := strings.Title(parts[0])
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
		SecurityGroup: ecs.SecurityGroupRule{
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
		runes := []rune(string(entry.SecurityGroup.Policy))
		runes[0] = unicode.ToLower(runes[0])
		policy = string(runes)
	}
	var direction string = entry.SecurityGroup.Direction
	var ipProtocol string = strings.ToLower(entry.SecurityGroup.IpProtocol)
	var portRange string = entry.SecurityGroup.PortRange
	var directionWord string
	if entry.SecurityGroup.Direction == "ingress" {
		directionWord = "from"
	} else {
		directionWord = "to"
	}
	var cidrIp string = entry.SecurityGroup.CidrIp
	var priority string = entry.SecurityGroup.Priority
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
	if entry.SecurityGroup.Description != "" {
		str += " # " + entry.SecurityGroup.Description
	}

	return str
}
