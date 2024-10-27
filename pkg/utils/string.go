package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// Converts the first letter of string to lower case
func LowercaseFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// Converts the first letter of string to upper case
func UppercaseFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// String stores v in a new string value and returns a pointer to it.
func String(v string) *string {
	return &v
}

// CamelCase converts string to camel case.
// Credits: This function is a slightly modified version of CamelCase function of the "github.com/samber/lo" package
func CamelCase(str string) string {
	if len(str) == 0 || str[0] == '_' {
		return str
	}
	items := Words(str)
	for i, item := range items {
		item = strings.ToLower(item)
		if i > 0 {
			item = UppercaseFirst(item)
		}
		items[i] = item
	}
	return strings.Join(items, "")
}

var (
	splitWordReg         = regexp.MustCompile(`([a-z])([A-Z0-9])|([a-zA-Z])([0-9])|([0-9])([a-zA-Z])|([A-Z])([A-Z])([a-z])`)
	splitNumberLetterReg = regexp.MustCompile(`([0-9])([a-zA-Z])`)
)

// Words splits string into an array of its words.
func Words(str string) []string {
	str = splitWordReg.ReplaceAllString(str, `$1$3$5$7 $2$4$6$8$9`)
	// example: Int8Value => Int 8Value => Int 8 Value
	str = splitNumberLetterReg.ReplaceAllString(str, "$1 $2")
	var result strings.Builder
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}
	return strings.Fields(result.String())
}
