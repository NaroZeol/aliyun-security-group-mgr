package reloader

import (
	"aliyun-security-group-mgr/internal/ecs"
	"os"
	"strings"
	"testing"
	"time"
)

var testGroup = []struct {
	line  string
	entry Entry
}{
	{
		line: "accept ingress tcp 22/22 from 0.0.0.0/0 priority 100 until 2024-12-31T23:59:59+08:00 # SSH access",
		entry: Entry{
			SecurityGroup: ecs.SecurityGroupRule{
				Policy:      ecs.PolicyAccept,
				Direction:   "ingress",
				IpProtocol:  "TCP",
				PortRange:   "22/22",
				CidrIp:      "0.0.0.0/0",
				Priority:    "100",
				Description: "SSH access",
			},
			ExpireAt: time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local),
		},
	},
	{
		line: "drop egress udp 53/53 to 0.0.0.0/0 priority 200 until 2024-12-31T23:59:59+08:00",
		entry: Entry{
			SecurityGroup: ecs.SecurityGroupRule{
				Policy:      ecs.PolicyDrop,
				Direction:   "egress",
				IpProtocol:  "UDP",
				PortRange:   "53/53",
				CidrIp:      "0.0.0.0/0",
				Priority:    "200",
				Description: "",
			},
			ExpireAt: time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local),
		},
	},
	{
		line: "accept ingress tcp 80/80 from 1.2.3.4/10 priority 300 until 2024-12-31T23:59:59+08:00 # TEST access",
		entry: Entry{
			SecurityGroup: ecs.SecurityGroupRule{
				Policy:      ecs.PolicyAccept,
				Direction:   "ingress",
				IpProtocol:  "TCP",
				PortRange:   "80/80",
				CidrIp:      "1.2.3.4/10",
				Priority:    "300",
				Description: "TEST access",
			},
			ExpireAt: time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local),
		},
	},
}

func TestDecodeEntry(t *testing.T) {
	for _, test := range testGroup {
		entry, err := DecodeEntry(test.line)
		if err != nil {
			t.Errorf("DecodeEntry(%q) returned error: %v", test.line, err)
			continue
		}

		if entry.SecurityGroup.Policy != test.entry.SecurityGroup.Policy ||
			entry.SecurityGroup.Direction != test.entry.SecurityGroup.Direction ||
			entry.SecurityGroup.IpProtocol != test.entry.SecurityGroup.IpProtocol ||
			entry.SecurityGroup.PortRange != test.entry.SecurityGroup.PortRange ||
			entry.SecurityGroup.CidrIp != test.entry.SecurityGroup.CidrIp ||
			entry.SecurityGroup.Priority != test.entry.SecurityGroup.Priority ||
			// !entry.ExpireAt.Equal(test.entry.ExpireAt) || // Timezones might be tricky, checking string rep if fails
			!entry.ExpireAt.Equal(test.entry.ExpireAt) ||
			entry.SecurityGroup.Description != test.entry.SecurityGroup.Description {
			t.Errorf("DecodeEntry(%q) = %+v; want %+v", test.line, entry, test.entry)
		}
	}
}

func TestEncodeEntry(t *testing.T) {
	for _, test := range testGroup {
		line := EncodeEntry(test.entry)
		if line != test.line {
			t.Errorf("EncodeEntry(%+v) = %q; want %q", test.entry, line, test.line)
		}
	}
}

func TestReadEntriesFromFile(t *testing.T) {
	// Create temp file
	tmpfile, err := os.CreateTemp("", "entries")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write content
	var content strings.Builder
	for _, test := range testGroup {
		content.WriteString(test.line + "\n")
	}
	if _, err := tmpfile.Write([]byte(content.String())); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Read
	entries, err := ReadEntriesFromFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ReadEntriesFromFile() error = %v", err)
	}

	if len(entries) != len(testGroup) {
		t.Errorf("ReadEntriesFromFile() returned %d entries, want %d", len(entries), len(testGroup))
	}
}

func TestEntry_EqualContent(t *testing.T) {
	e1 := Entry{
		SecurityGroup: ecs.SecurityGroupRule{
			Policy: "Accept",
			CidrIp: "0.0.0.0/0",
		},
		ExpireAt: time.Now(),
	}
	e2 := Entry{
		SecurityGroup: ecs.SecurityGroupRule{
			Policy: "Accept",
			CidrIp: "0.0.0.0/0",
		},
		ExpireAt: time.Now().Add(1 * time.Hour), // Different time
	}
	e3 := Entry{
		SecurityGroup: ecs.SecurityGroupRule{
			Policy: "Drop",
			CidrIp: "0.0.0.0/0",
		},
		ExpireAt: time.Now(),
	}

	if !e1.EqualContent(e2) {
		t.Error("e1 should be equal content to e2 (ignoring time)")
	}
	if e1.EqualContent(e3) {
		t.Error("e1 should not be equal content to e3")
	}
}
