package main

import (
	"fmt"
	"os"

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

func read() (string, ReadResult) {
	os.Stdout.Write([]byte{'$', ' '})

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

		case '\t':
			result := autocomplete(&line)
			switch result {
			case AutoCompleteNone:
				break
			case AutoCompleteFound:
				break
			case AutoCompleteMore:
				break
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
	arguments, redirection := parser.Parse()
	if len(arguments) <= 0 {
		return
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
