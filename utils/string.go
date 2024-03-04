package utils

import (
	"strings"
	"unicode"
)

func FirstToUpper(text string) string {
	if text == "" {
		return text
	}
	return strings.ToUpper(text[:1]) + text[1:]
}

func FirstToLower(text string) string {
	if text == "" {
		return text
	}
	return strings.ToLower(text[:1]) + text[1:]
}

func IsFirstUpper(str string) bool {
	if len(str) == 0 {
		return false
	}
	first := str[0]
	return unicode.IsUpper(rune(first))
}

func IsFirstLower(str string) bool {
	return !IsFirstUpper(str)
}
