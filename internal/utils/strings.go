package utils

import (
	"strings"
)

// extract comments from a line
func ExtractCommentFromLine(line string) string {
	if idx := strings.Index(line, "#"); idx != -1 {
		return strings.TrimSpace(line[idx+1:])
	}
	return ""
}

// remove comments from a line
func RemoveCommentFromLine(line string) string {
	if idx := strings.Index(line, "#"); idx != -1 {
		line = line[:idx]
	}
	return strings.TrimSpace(line)
}
