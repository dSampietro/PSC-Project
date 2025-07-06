package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func filter[T any](s []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(s)) // Pre-allocate for efficiency
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func isWhitespace(ch string) bool {
	runes := []rune(ch)
	if len(runes) == 0 {
		return false
	}
	return unicode.IsSpace(runes[0])
}

//Tokenize a string into its non-whitespace components
func Tokenize(s string) []string {
	substrings := strings.Split(s, "")
	
	return filter(substrings, func(ch string) bool {
		return !isWhitespace(ch)
	})
}

func ParseFile(filename string) string {
	text, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Cannot succesfully read file")
	}

	return string(text)
}
