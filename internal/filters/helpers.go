package filters

import (
	"strconv"
	"strings"
)

// parseBitrateValue converts a bitrate string (e.g., "2000kbps", "2mbps") to bits per second
func parseBitrateValue(valueStr string) int64 {
	valueStr = strings.ToLower(strings.TrimSpace(valueStr))
	valueStr = strings.ReplaceAll(valueStr, " ", "")

	if strings.HasSuffix(valueStr, "kbps") || strings.HasSuffix(valueStr, "kb/s") {
		numStr := strings.TrimSuffix(strings.TrimSuffix(valueStr, "kbps"), "kb/s")
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0
		}
		return int64(num * 1000)
	} else if strings.HasSuffix(valueStr, "mbps") || strings.HasSuffix(valueStr, "mb/s") {
		numStr := strings.TrimSuffix(strings.TrimSuffix(valueStr, "mbps"), "mb/s")
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0
		}
		return int64(num * 1000000)
	} else {
		num, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return 0
		}
		return num
	}
}

// compareNumeric compares two numeric values using the specified operator
func compareNumeric(actual int64, operator string, target int64) bool {
	switch operator {
	case "IS", "==":
		return actual == target
	case "IS_NOT", "!=":
		return actual != target
	case ">":
		return actual > target
	case ">=":
		return actual >= target
	case "<":
		return actual < target
	case "<=":
		return actual <= target
	default:
		return false
	}
}

// compareFloat compares two float values using the specified operator
func compareFloat(actual float64, operator string, target float64) bool {
	switch operator {
	case "IS", "==":
		return actual == target
	case "IS_NOT", "!=":
		return actual != target
	case ">":
		return actual > target
	case ">=":
		return actual >= target
	case "<":
		return actual < target
	case "<=":
		return actual <= target
	default:
		return false
	}
}

// compareString compares two strings using the specified operator
func compareString(actual string, operator string, target string) bool {
	actual = strings.ToLower(strings.TrimSpace(actual))
	target = strings.ToLower(strings.TrimSpace(target))

	switch operator {
	case "IS", "==":
		return actual == target
	case "IS_NOT", "!=":
		return actual != target
	case "CONTAINS":
		return strings.Contains(actual, target)
	case "NOT_CONTAINS":
		return !strings.Contains(actual, target)
	default:
		return false
	}
}

// compareBool compares two boolean values using the specified operator
func compareBool(actual bool, operator string, valueStr string) bool {
	target := strings.ToLower(valueStr) == "true" || valueStr == "1"

	switch operator {
	case "IS", "==":
		return actual == target
	case "IS_NOT", "!=":
		return actual != target
	default:
		return false
	}
}
