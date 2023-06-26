package land

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func latinize(value string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, value)
	if err != nil {
		return ""
	}
	return result
}

func replaceSpecialCharacters(value string) string {
	replacer := regexp.MustCompile(`[-_.,=&;@/(){}]`)
	marksReplacer := strings.NewReplacer("'", "’")
	value = replacer.ReplaceAllString(value, " ")
	value = marksReplacer.Replace(value)
	return value
}

func simplify(value string) string {
	value = strings.ToLower(value)
	value = latinize(value)
	return value
}
