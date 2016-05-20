package config

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var intervalRegex = regexp.MustCompile(`([0-9]+)([smhdw])`)

// ParseTimeInterval replaces the builtin parser by adding support for days and weeks
func ParseTimeInterval(source string) (time.Duration, error) {
	matches := intervalRegex.FindStringSubmatch(strings.ToLower(source))

	if matches == nil {
		return 0, errors.New("Invalid time interval " + source)
	}

	val, err := strconv.Atoi(matches[1])

	if err != nil {
		return 0, err
	}

	switch matches[2] {
	case "s":
		return time.Duration(val) * time.Second, nil

	case "m":
		return time.Duration(val) * time.Second * 60, nil

	case "h":
		return time.Duration(val) * time.Second * 60 * 60, nil

	case "d":
		return time.Duration(val) * time.Second * 60 * 60 * 24, nil

	case "w":
		return time.Duration(val) * time.Second * 60 * 60 * 24 * 7, nil

	default:
		return 0, errors.New("Invalid time interval " + source)
	}
}
