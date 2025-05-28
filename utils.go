package telegram

import (
	"regexp"
	"strconv"
	"strings"
)

// IsFloodError checks if error is a flood error
func IsFloodError(err error) bool {
	return strings.Contains(err.Error(), "Too Many Requests")
}

// ParseRetryAfter returns retry after seconds
func ParseRetryAfter(err error) int {
	re := regexp.MustCompile(`retry after (\d+)`)
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) == 2 {
		sec, _ := strconv.Atoi(matches[1])
		return sec
	}
	return 10 // fallback
}
