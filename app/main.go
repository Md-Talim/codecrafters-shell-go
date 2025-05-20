package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/md-talim/codecrafters-shell-go/app/parser"
	"github.com/md-talim/codecrafters-shell-go/app/shellio"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

type ReadResult int

const (
	ReadResultQuit ReadResult = iota
	ReadResultEmpty
	ReadResultContent
)

func prompt() {
	os.Stdout.Write([]byte{'$', ' '})
}

func read() (string, ReadResult) {
	prompt()

	var stdinFd = os.Stdin.Fd()
	var previous unix.Termios
	if err := termios.Tcgetattr(stdinFd, &previous); err != nil {
		panic(err)
	}

	var new = unix.Termios(previous)
	new.Iflag &= unix.IGNCR  // ignore recieved CR
	new.Lflag ^= unix.ICANON // disable canonical mode
	new.Lflag ^= unix.ECHO   // disable echo of input
	new.Cc[unix.VMIN] = 1
	new.Cc[unix.VTIME] = 0
	if err := termios.Tcsetattr(stdinFd, termios.TCSANOW, &new); err != nil {
		panic(err)
	}
	defer termios.Tcsetattr(stdinFd, termios.TCSANOW, &previous)

	line := ""
	bellRang := false
	buffer := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(buffer)
		if err != nil {
			return "", ReadResultQuit
		}

		character := buffer[0]

		switch character {
		case 0x4: // CTRL + D
			return "", ReadResultQuit
		case '\r': // ENTER
			fallthrough
		case '\n': // NEW LINE
			os.Stdout.Write([]byte{'\r', '\n'})
			if len(line) == 0 {
				return "", ReadResultEmpty
			} else {
				return line, ReadResultContent
			}

		case '\t': // TAB
			result := autocomplete(&line, bellRang)
			switch result {
			case AutoCompleteNone:
				bellRang = false
				bell()
			case AutoCompleteFound:
				bellRang = false
			case AutoCompleteMore:
				bellRang = true
				bell()
			}

		case 0x1b: // ARROW KEYS
			os.Stdin.Read(buffer) // '['
			os.Stdin.Read(buffer) // 'A' or 'B' or 'C' or 'D'

		case 0x7f: // BACKSPACE
			if len(line) != 0 {
				line = line[:len(line)-1]
				os.Stdout.Write([]byte{'\b', ' ', '\b'})
			}
		default:
			os.Stdout.Write(buffer)
			line += string(character)
		}
	}
}

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

func executePipeCommand(parsedCommands [][]string, shellio shellio.IO) {
	command1Def := parsedCommands[0]
	command2Def := parsedCommands[1]

	if len(command1Def) == 0 || len(command2Def) == 0 {
		fmt.Fprintln(shellio.ErrorFile(), "shell: error, empty command in pipeline")
		return
	}

	command1Name, command1Args := command1Def[0], command1Def[1:]
	command2Name, command2Args := command2Def[0], command2Def[1:]

	// Create the pipe
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(shellio.ErrorFile(), "shell: error creating pipe: ", err)
		return
	}

	// First command
	cmd1 := exec.Command(command1Name, command1Args...)
	cmd1.Stdout = w
	cmd1.Stderr = shellio.ErrorFile()

	// Second command
	cmd2 := exec.Command(command2Name, command2Args...)
	cmd2.Stdin = r
	cmd2.Stdout = shellio.OutputFile()
	cmd2.Stderr = shellio.ErrorFile()

	if err := cmd1.Start(); err != nil {
		fmt.Fprintf(shellio.ErrorFile(), "shell: error starting command %s: %v\n", command1Name, err)
		w.Close()
		r.Close()
		return
	}
	// Close the write end of the pipe in the parent. cmd1 still holds it open.
	// This is crucial so that cmd2 receives EOF when cmd1 finishes.
	w.Close()

	if err := cmd2.Start(); err != nil {
		fmt.Fprintf(shellio.ErrorFile(), "shell: error starting command %s: %v\n", command2Name, err)
		r.Close()   // Close read end as well
		cmd1.Wait() // Attempt to reap cmd1 if it started
		return
	}

	// Close the read end of the pipe in the parent. cmd2 still holds it open.
	r.Close()

	// Wait for both commands to finish
	if err := cmd1.Wait(); err != nil {
		fmt.Fprintf(shellio.ErrorFile(), "shell: error waiting for command %s: %v\n", command1Name, err)
	}
	if err := cmd2.Wait(); err != nil {
		fmt.Fprintf(shellio.ErrorFile(), "shell: error waiting for command %s: %v\n", command2Name, err)
	}
}

func executeSingleCommand(parsedCommands [][]string, shellio shellio.IO) {
	args := parsedCommands[0]
	if len(args) == 0 {
		return
	}
	commandName := args[0]
	commandArgs := args[1:]

	if builtinCommand, ok := builtinCommands[commandName]; ok {
		builtinCommand(commandArgs, shellio)
	} else if _, ok := findPath(commandName); ok {
		externelCommand(commandName, commandArgs, shellio)
	} else {
		fmt.Fprintf(shellio.OutputFile(), "%s: command not found\n", commandName)
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
