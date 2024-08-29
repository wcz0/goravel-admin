package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	// 处理 JSON 字符串
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type for StringSlice")
	}

	var strSlice []string
	if err := json.Unmarshal(bytes, &strSlice); err != nil {
		return err
	}
	*s = strSlice
	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}