package land

import (
	"fmt"
	"slices"
	"strings"
)

func createTSQuery(value string) string {
	result := make([]string, 0)
	for _, item := range strings.Split(replaceSpecialCharacters(value), " ") {
		if len(item) == 0 {
			continue
		}
		item = simplify(item)
		if !slices.Contains(result, item) {
			result = append(result, item+":*")
		}
	}
	return fmt.Sprintf("to_tsquery('%s')", strings.Join(result, " & "))
}

func createTSVectors(values ...any) string {
	result := make([]string, 0)
	for _, v := range values {
		s := fmt.Sprintf("%v", v)
		s = replaceSpecialCharacters(s)
		s = simplify(s)
		if !slices.Contains(result, s) && v != nil {
			result = append(result, s)
		}
	}
	return fmt.Sprintf("to_tsvector('%s')", strings.Join(result, " "))
}
