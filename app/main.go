package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
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
		case "pwd":
			handlePwdCommad()
		default:
			if _, err := findPath(command); err == nil {
				handleExternelCommand(command, arguments)
			} else {
				fmt.Printf("%s: command not found\n", command)
			}
		}
	}
}

func findPath(command string) (string, error) {
	PATH := os.Getenv("PATH")
	directories := strings.SplitSeq(PATH, string(os.PathListSeparator))

	for dir := range directories {
		fullPath := path.Join(dir, command)
		if fileInfo, err := os.Stat(fullPath); err == nil && fileInfo.Mode().IsRegular() && (fileInfo.Mode().Perm()&0111 != 0) {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("%s: not found", command)
}

func handleExternelCommand(command string, args []string) {
	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func handleTypeCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("type: missing operand")
		return
	}

	for _, arg := range args {
		switch arg {
		case "exit", "echo", "type", "pwd":
			fmt.Printf("%s is a shell builtin\n", arg)
		default:
			if path, err := findPath(arg); err == nil {
				fmt.Printf("%s is %s\n", arg, path)
			} else {
				fmt.Println(err)
			}
		}
	}
}

func handlePwdCommad() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(dir)
}
