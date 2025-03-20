package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
        fmt.Fprintln(os.Stderr, "Error reading input: ", err)
		os.Exit(0)
	}

	// Command not found
	fmt.Printf("%s: command not found", command[:len(command)-1])
}
