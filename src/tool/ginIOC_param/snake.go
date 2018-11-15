package ginIOC_param

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func toSnake(name string, do bool) string {
	if !do {
		return name
	}
	if name == "" {
		return ""
	}
	if unicode.IsUpper(rune(name[0])) {
		return name
	}

	var i, j, p int
	var words []string

	for i < len(name) {
		if name[i] != '_' {
			break
		}
		i++
		j++
		words = append(words, "")
	}

	for p = i; p < len(name); p++ {
		if name[p] == '_' {
			words = append(words, strings.ToLower(name[j:p]))
			j = p + 1
			i = p + 1
			continue
		}
		if unicode.IsUpper(rune(name[p])) {
			if j != i {
				words = append(words, strings.ToLower(name[j:p]))
				j = p
				i = p
			}
		} else {
			i = p
		}
	}
	if j < p {
		words = append(words, strings.ToLower(name[j:p]))
	}
	return strings.Join(words, "_")
}
