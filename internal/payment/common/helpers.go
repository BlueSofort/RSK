package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ReadString 从 map 中读取字符串值。
func ReadString(raw map[string]interface{}, key string) string {
	if raw == nil || key == "" {
		return ""
	}
	value, ok := raw[key]
	if !ok || value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ReadInt64 从 map 中读取 int64 值。
func ReadInt64(raw map[string]interface{}, key string) (int64, bool) {
	if raw == nil || key == "" {
		return 0, false
	}
	value, ok := raw[key]
	if !ok || value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case float64:
		return int64(v), true
	case int64:
		return v, true
	case json.Number:
		parsed, err := v.Int64()
		if err != nil {
			return 0, false
		}
		return parsed, true
	case string:
		parsed, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

// ReadMap 从 map 中读取子 map。
func ReadMap(raw map[string]interface{}, key string) map[string]interface{} {
	if raw == nil || key == "" {
		return nil
	}
	value, ok := raw[key]
	if !ok || value == nil {
		return nil
	}
	if m, ok := value.(map[string]interface{}); ok {
		return m
	}
	return nil
}
