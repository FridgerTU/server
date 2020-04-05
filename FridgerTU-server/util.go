package main

import "strings"

func containsInvalidArgs(query []string) bool {
	for _, arg := range query {
		if strings.HasPrefix(arg, ",") {
			return true
		}
	}
	return false
}

func containsEmptyArgs(query []string) bool {
	for _, arg := range query {
		if arg == "" {
			return true
		}
	}
	return false
}
