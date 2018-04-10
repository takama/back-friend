package helper

import (
	"regexp"
	"strings"
)

// ToSnake converts string to underscore/snake string
func ToSnake(str string) string {
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").
		ReplaceAllString(regexp.MustCompile("(.)([A-Z][a-z]+)").
			ReplaceAllString(str, "${1}_${2}"), "${1}_${2}"))
}
