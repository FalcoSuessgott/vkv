package regex

import (
	"regexp"
)

// MatchRegex matches a given string to a regex.
// https://pkg.go.dev/regexp/syntax
func MatchRegex(pattern, s string) (bool, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	return r.MatchString(s), nil
}
