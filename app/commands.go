package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type BuiltinCommand func([]string)

var builtinCommands map[string]BuiltinCommand

func exitCommand([]string) {
	os.Exit(0)
}

func echoCommand(args []string) {
	output := strings.Join(args, " ")
	fmt.Println(output)
}

func typeCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("type: missing operand")
		return
	}

	for _, arg := range args {
		switch arg {
		case "exit", "echo", "type", "pwd":
			fmt.Printf("%s is a shell builtin\n", arg)
		default:
			if path, ok := findPath(arg); ok {
				fmt.Printf("%s is %s\n", arg, path)
			} else {
				fmt.Printf("%s: not found\n", arg)
			}
		}
	}
}

func pwdCommand(_ []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(dir)
}

func externelCommand(command string, args []string) {
	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}

}

func findPath(command string) (string, bool) {
	PATH := os.Getenv("PATH")
	directories := strings.SplitSeq(PATH, string(os.PathListSeparator))

	for dir := range directories {
		fullPath := path.Join(dir, command)
		if fileInfo, err := os.Stat(fullPath); err == nil && fileInfo.Mode().IsRegular() && (fileInfo.Mode().Perm()&0111 != 0) {
			return fullPath, true
		}
	}
	return "", false
}
