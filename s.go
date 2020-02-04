package gometrics

import (
	"bytes"
	"encoding/json"
	"strings"
)

// StripAny strips any Unicode code points in chars are within s.
func StripAny(s, chars string) string {
	filter := func(r rune) rune {
		if !strings.ContainsRune(chars, r) {
			return r
		}

		return -1
	}

	return strings.Map(filter, s)
}

// Esc escapes s to a human readable format
func Esc(s string) string {
	j, _ := json.Marshal(s)
	return string(j)[1 : len(j)-1]
}

// Abbr abbreviate s to max length
func Abbr(s string, max int, postfix string) string {
	l := len(s)
	if l <= max {
		return s
	}

	i := max - len(postfix)
	if i > 0 {
		return s[0:i] + postfix
	}

	return postfix
}

// JSONCompact compact the JSON encoding of data silently
func JSONCompact(data interface{}) string {
	return PickFirst(JSONCompactE(data))
}

// JSONCompactE compact the JSON encoding of data
func JSONCompactE(data interface{}) (string, error) {
	switch v := data.(type) {
	case string:
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, []byte(v)); err != nil {
			return "", err
		}

		return buffer.String(), nil
	case []byte:
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, v); err != nil {
			return "", err
		}

		return buffer.String(), nil
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return string(b), nil
	}
}

// PickFirst ignores the error and returns s
func PickFirst(s string, _ interface{}) string {
	return s
}
