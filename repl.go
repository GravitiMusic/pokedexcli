package main

import (
	"strings"
)

func cleanInput(text string) []string {
	s := strings.TrimSpace(text)
	lower_s := strings.ToLower(s)
	words := strings.Split(lower_s, " ")
	return words
}
