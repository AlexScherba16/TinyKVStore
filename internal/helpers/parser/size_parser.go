package parser

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	kb = "kb"
	mb = "mb"
)

func ParseBufferSize(size string) int {
	re := regexp.MustCompile(`^\s*(\d+)\s*([a-zA-Z]+)\s*$`)
	matches := re.FindAllStringSubmatch(size, -1)
	if matches == nil {
		return -1
	}

	if len(matches) > 1 || len(matches[0]) != 3 {
		return -1
	}

	sizeValue, err := strconv.ParseInt(matches[0][1], 10, 32)
	if err != nil {
		return -1
	}

	multiplier := strings.ToLower(strings.TrimSpace(matches[0][2]))
	switch multiplier {
	case kb:
		return int(sizeValue) * 1024
	case mb:
		return int(sizeValue) * 1024 * 1024
	}

	return -1
}
