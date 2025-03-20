package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			os.Exit(0)
		}

		input = input[:len(input)-1]           // Remove the '\n' at the end
		arguments := strings.Split(input, " ") // Split input by spaces

		command := arguments[0]
		arguments = arguments[1:]
		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			output := strings.Join(arguments, " ")
			fmt.Println(output)
		default: // Command not found
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
