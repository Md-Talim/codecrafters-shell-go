package executor

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/internal/shellio"
	"github.com/md-talim/codecrafters-shell-go/internal/utils"
)

type BuiltinCommandExecutor func([]string, shellio.IO)
type BuiltinCommandsMap map[string]BuiltinCommandExecutor

var builtinCommands BuiltinCommandsMap
var shellHistory []string

func init() {
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
	return len(shellHistory)
}

func GetHistoryEntry(index int) (string, bool) {
	if index >= 0 && index < len(shellHistory) {
		return shellHistory[index], true
	}
	return "", false
}

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
		case "exit", "echo", "type", "pwd", "cd", "history":
			fmt.Fprintf(io.OutputFile(), "%s is a shell builtin\n", arg)
		default:
			if path, ok := utils.FindPath(arg); ok {
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
	limit, err := getHistoryLimit(args)
	if err != nil {
		fmt.Fprintln(io.ErrorFile(), err)
		return
	}

	startIndex := len(shellHistory) - limit
	for i := startIndex; i < len(shellHistory); i++ {
		fmt.Fprintf(io.OutputFile(), "    %d  %s\n", i+1, shellHistory[i])
	}
}
