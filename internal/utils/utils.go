package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func InArray(needle string, array []string) bool {
	for _, str := range array {
		if needle == str {
			return true
		}
	}
	return false
}

func SplitSemVer(version string) (int, int, int, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("incorrect format: %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid value: %v", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid value: %v", err)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid value: %v", err)
	}

	return major, minor, patch, nil
}
