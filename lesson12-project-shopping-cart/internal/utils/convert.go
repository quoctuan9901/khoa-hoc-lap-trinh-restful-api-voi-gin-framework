package utils

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)[A-Z][a-z]+")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func CamelToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func NormalizeString(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

func ConvertToInt32Pointer(value int) *int32 {
	if value == 0 {
		return nil
	}

	v := int32(value)
	
	return &v
}
