package main

import (
	"github.com/md-talim/codecrafters-shell-go/app/parser"
	"github.com/md-talim/codecrafters-shell-go/app/shellio"
)

type ReadResult int

const (
	ReadResultQuit ReadResult = iota
	ReadResultEmpty
	ReadResultContent
)

func eval(input string) {
	parser := parser.NewParser(input)
	parsedCommands, redirectionConfig := parser.Parse()
	if len(parsedCommands) <= 0 {
		return
	}

	shellio, isRedirected := shellio.OpenIo(redirectionConfig)
	if isRedirected {
		defer shellio.Close()
	}

	if len(parsedCommands) == 1 {
		executeSingleCommand(parsedCommands, shellio)
	} else if len(parsedCommands) == 2 {
		executePipeCommand(parsedCommands, shellio)
	}
}

func main() {
	builtinCommands = make(map[string]BuiltinCommand)
	builtinCommands["exit"] = exitCommand
	builtinCommands["echo"] = echoCommand
	builtinCommands["type"] = typeCommand
	builtinCommands["pwd"] = pwdCommand
	builtinCommands["cd"] = cdCommand

	for {
		input, result := read()

		switch result {
		case ReadResultQuit:
			return
		case ReadResultEmpty:
			continue
		case ReadResultContent:
			eval(input)
		}
	}
}
