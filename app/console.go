package main

import (
	"fmt"
	"os"

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
		fmt.Fprintln(os.Stderr, "Error initializing terminal: ", err)
		return "", ReadResultQuit
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
