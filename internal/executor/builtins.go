package executor

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/internal/shellio"
)

type BuiltinCommandExecutor func([]string, shellio.IO)
type BuiltinCommandsMap map[string]BuiltinCommandExecutor

var builtinCommands BuiltinCommandsMap
var history CommandHistory

func init() {
	loadHistoryFromHISTFILE()
	builtinCommands = BuiltinCommandsMap{
		"cd":      cdCommand,
		"echo":    echoCommand,
		"exit":    exitCommand,
		"history": historyCommand,
		"pwd":     pwdCommand,
		"type":    typeCommand,
	}
}

func BuiltinCommands() BuiltinCommandsMap {
	return builtinCommands
}

func GetHistoryLength() int {
	return history.length()
}

func GetHistoryEntry(index int) (string, bool) {
	if index >= 0 && index < history.length() {
		return history.at(index), true
	}
	return "", false
}

func exitCommand(_ []string, _ shellio.IO) {
	writeHistoryToHISTFILE()
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
		case "exit", "echo", "type", "pwd", "cd", "history":
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
			fmt.Fprintln(io.ErrorFile(), "cd: $HOME is not set.")
		} else {
			newDir = path.Join(HOME, newDir[1:])
		}
	}

	if err := os.Chdir(newDir); err != nil {
		fmt.Fprintf(io.ErrorFile(), "cd: %s: No such file or directory\n", newDir)
	}
}

func historyCommand(args []string, io shellio.IO) {
	if len(args) == 0 {
		history.printAll(io)
		return
	}

	// The first arg can be the action like "-r", "-w", or "-a"
	// It can also be the limit for history
	action := args[0]
	if action == "-r" {
		history.appendFromFile(args[1])
	} else if action == "-w" {
		history.saveToFile(args[1])
	} else if action == "-a" {
		history.appendToFile(args[1])
	} else {
		limit, err := parseHistoryLimit(action)
		if err != nil {
			fmt.Fprintln(io.ErrorFile(), err)
			return
		}
		history.printLast(limit, io)
	}
}
