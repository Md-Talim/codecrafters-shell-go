package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/app/parser"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			os.Exit(0)
		}

		input = input[:len(input)-1] // Remove the '\n' at the end
		parser := parser.NewParser(input)
		arguments := parser.Parse()
		if len(arguments) <= 0 {
			continue
		}

		command := arguments[0]
		arguments = arguments[1:]
		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			output := strings.Join(arguments, " ")
			fmt.Println(output)
		case "type":
			handleTypeCommand(arguments)
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func handleTypeCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("type: missing operand")
		return
	}

	for _, arg := range args {
		switch arg {
		case "exit", "echo", "type":
			fmt.Printf("%s is a shell builtin\n", arg)
		default:
			fmt.Printf("%s: not found\n", arg)
		}
	}
}
