package main

import "strings"

// RemoveURLs removes all URLs from the input string
func RemoveURLs(descRaw string) string {
	words := strings.Fields(descRaw)
	var filtered []string
	for _, word := range words {
		if !(strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://")) {
			filtered = append(filtered, word)
		}
	}
	return strings.Join(filtered, " ")
}
