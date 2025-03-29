package main

import (
	"os"
	"strings"
)

type AutoCompleteResult int

const (
	AutoCompleteNone AutoCompleteResult = iota
	AutoCompleteFound
	AutoCompleteMore
)

func printCompletion(line *string, completion string) {
	os.Stdout.WriteString(completion)
	*line += completion

	os.Stdout.WriteString(" ")
	*line += " "
}

func autocomplete(line *string) AutoCompleteResult {
	var completions []string

	for name := range builtinCommands {
		if strings.HasPrefix(name, *line) {
			completion := name[len(*line):]
			completions = append(completions, completion)
		}
	}

	if len(completions) == 0 {
		return AutoCompleteNone
	}

	if len(completions) == 1 {
		completion := completions[0]
		printCompletion(line, completion)
		return AutoCompleteNone
	}

	return AutoCompleteNone
}
