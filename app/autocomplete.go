package main

import (
	"os"
	"slices"
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

	PATH := os.Getenv("PATH")
	directories := strings.SplitSeq(PATH, string(os.PathListSeparator))

	for directory := range directories {
		entries, err := os.ReadDir(directory)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasPrefix(name, *line) {
				continue
			}

			path := strings.Join([]string{directory, name}, "/")
			stat, err := os.Stat(path)
			if err != nil || !stat.Mode().IsRegular() || stat.Mode().Perm()&0111 == 0 {
				continue
			}

			completion := name[len(*line):]
			if !slices.Contains(completions, completion) {
				completions = append(completions, completion)
			}
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

func bell() {
	os.Stdout.Write([]byte{'\a'})
}
