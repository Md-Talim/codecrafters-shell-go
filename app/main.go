package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/md-talim/codecrafters-shell-go/app/parser"
	"github.com/md-talim/codecrafters-shell-go/app/shellio"
)

func main() {
	builtinCommands = make(map[string]BuiltinCommand)
	builtinCommands["exit"] = exitCommand
	builtinCommands["echo"] = echoCommand
	builtinCommands["type"] = typeCommand
	builtinCommands["pwd"] = pwdCommand
	builtinCommands["cd"] = cdCommand

	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			os.Exit(0)
		}

		input = input[:len(input)-1] // Remove the '\n' at the end
		parser := parser.NewParser(input)
		arguments, redirection := parser.Parse()
		if len(arguments) <= 0 {
			continue
		}

		shellio, ok := shellio.OpenIo(redirection)
		if ok {
			defer shellio.Close()
		}

		command := arguments[0]
		arguments = arguments[1:]

		if builtinCommand, ok := builtinCommands[command]; ok {
			builtinCommand(arguments, shellio)
		} else if _, ok := findPath(command); ok {
			externelCommand(command, arguments, shellio)
		} else {
			fmt.Fprintf(shellio.OutputFile(), "%s: command not found\n", command)
		}
	}
}
