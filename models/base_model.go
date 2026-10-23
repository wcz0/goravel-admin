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

    // 支持 []byte / string 两种来源
    var raw string
    switch v := value.(type) {
    case []byte:
        raw = string(v)
    case string:
        raw = v
    default:
        return errors.New("invalid type for StringSlice")
    }

    // 优先尝试标准 JSON []string
    var strSlice []string
    if err := json.Unmarshal([]byte(raw), &strSlice); err == nil {
        *s = strSlice
        return nil
    }

    // 尝试 JSON 字符串包裹的数组（双重编码）
    var inner string
    if err := json.Unmarshal([]byte(raw), &inner); err == nil {
        // 以 inner 继续解析
        if err := json.Unmarshal([]byte(inner), &strSlice); err == nil {
            *s = strSlice
            return nil
        }
        // 单引号数组
        normalized := replaceSingleQuotes(inner)
        if err := json.Unmarshal([]byte(normalized), &strSlice); err == nil {
            *s = strSlice
            return nil
        }
        parts := splitComma(inner)
        if len(parts) > 0 {
            *s = parts
            return nil
        }
    }

    // 兼容 PHP 单引号数组: ['GET','POST']
    if len(raw) > 0 && raw[0] == '[' {
        normalized := raw
        normalized = replaceSingleQuotes(normalized)
        if err := json.Unmarshal([]byte(normalized), &strSlice); err == nil {
            *s = strSlice
            return nil
        }
    }

    // 兼容逗号分隔: GET,POST
    parts := splitComma(raw)
    if len(parts) > 0 {
        *s = parts
        return nil
    }
    return errors.New("invalid StringSlice content")
}

func (s StringSlice) Value() (driver.Value, error) {
    return json.Marshal(s)
}

// replaceSingleQuotes 将单引号包裹的列表转换为双引号 JSON
func replaceSingleQuotes(s string) string {
    b := make([]byte, 0, len(s))
    for i := 0; i < len(s); i++ {
        if s[i] == '\'' { b = append(b, '"') } else { b = append(b, s[i]) }
    }
    return string(b)
}

func splitComma(s string) []string {
    var out []string
    cur := ""
    for i := 0; i < len(s); i++ {
        if s[i] == ',' {
            v := trimSpaceQuote(cur)
            if v != "" { out = append(out, v) }
            cur = ""
        } else { cur += string(s[i]) }
    }
    v := trimSpaceQuote(cur)
    if v != "" { out = append(out, v) }
    return out
}

func trimSpaceQuote(s string) string {
    // 简单去空格与引号
    for len(s) > 0 && (s[0] == ' ' || s[0] == '\'' || s[0] == '"') { s = s[1:] }
    for len(s) > 0 {
        c := s[len(s)-1]
        if c == ' ' || c == '\'' || c == '"' { s = s[:len(s)-1] } else { break }
    }
    return s
}