package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

type TagList []string

func (t TagList) Value() (driver.Value, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return strings.Join(t, ","), nil
}

func (t *TagList) Scan(value interface{}) error {
	if value == nil {
		*t = TagList{}
		return nil
	}

	var raw string
	switch v := value.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return errors.New("unsupported type for TagList")
	}

	if strings.TrimSpace(raw) == "" {
		*t = TagList{}
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make(TagList, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	*t = out
	return nil
}

func (t TagList) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(t))
}

func (t *TagList) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = TagList{}
		return nil
	}
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	out := make(TagList, 0, len(arr))
	for _, s := range arr {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	*t = out
	return nil
}
