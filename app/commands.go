package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/app/shellio"
)

type BuiltinCommand func([]string, shellio.IO)

var builtinCommands map[string]BuiltinCommand

func exitCommand(_ []string, _ shellio.IO) {
	os.Exit(0)
}

func echoCommand(args []string, io shellio.IO) {
	output := strings.Join(args, " ")
	fmt.Fprintln(io.OutputFile(), output)
}

func typeCommand(args []string, io shellio.IO) {
	if len(args) < 1 {
		fmt.Println("type: missing operand")
		return
	}

	for _, arg := range args {
		switch arg {
		case "exit", "echo", "type", "pwd":
			fmt.Fprintf(io.OutputFile(), "%s is a shell builtin\n", arg)
		default:
			if path, ok := findPath(arg); ok {
				fmt.Fprintf(io.OutputFile(), "%s is %s\n", arg, path)
			} else {
				fmt.Fprintf(io.OutputFile(), "%s: not found\n", arg)
			}
		}
	}
}

func pwdCommand(_ []string, io shellio.IO) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Fprintln(io.OutputFile(), dir)
}

func cdCommand(args []string, io shellio.IO) {
	newDir := args[0]

	if strings.HasPrefix(newDir, "~") {
		HOME := os.Getenv("HOME")
		if (len(HOME)) == 0 {
			fmt.Fprintln(io.OutputFile(), "cd: $HOME is not set.")
		} else {
			newDir = path.Join(HOME, newDir[1:])
		}
	}

	if err := os.Chdir(newDir); err != nil {
		fmt.Fprintf(io.OutputFile(), "cd: %s: No such file or directory\n", newDir)
	}
}

func externelCommand(command string, args []string, io shellio.IO) {
	cmd := exec.Command(command, args...)

	cmd.Stdout = io.OutputFile()
	cmd.Stderr = os.Stderr

	cmd.Run()
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
