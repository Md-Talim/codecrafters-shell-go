package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/internal/executor"
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

var historyNavigationIndex int

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
			historyNavigationIndex++
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
			// Read the next byte for '['
			bracketBuffer := make([]byte, 1)
			n, err := os.Stdin.Read(bracketBuffer)
			if err != nil || n == 0 || bracketBuffer[0] != '[' {
				continue
			}
			// Got 'ESC [', now read the arrow code
			arrowCodeBuffer := make([]byte, 1)
			n, err = os.Stdin.Read(arrowCodeBuffer)
			if err != nil || n == 0 {
				continue
			}

			switch arrowCodeBuffer[0] {
			case 'A': // Up Arrow
				if executor.GetHistoryLength() == 0 {
					bell()
					continue
				}
				if historyNavigationIndex > 0 {
					historyNavigationIndex--
					recalledCommand, ok := executor.GetHistoryEntry(historyNavigationIndex)
					if ok {
						currentVisualLength := len(line)
						fmt.Fprintf(os.Stdout, "\r%s\r", strings.Repeat(" ", len("$ ")+currentVisualLength))
						prompt()
						os.Stdout.WriteString(recalledCommand)
						line = recalledCommand
					} else {
						if historyNavigationIndex < executor.GetHistoryLength() {
							historyNavigationIndex++
						}
						bell()
					}
				} else {
					bell()
				}
			case 'B': // Down Arrow
				if executor.GetHistoryLength() == 0 || historyNavigationIndex <= 0 {
					bell()
					continue
				}
				historyNavigationIndex++
				recalledCommand, ok := executor.GetHistoryEntry(historyNavigationIndex)
				if ok {
					currentVisualLength := len(line)
					fmt.Fprintf(os.Stdout, "\r%s\r", strings.Repeat(" ", len("$ ")+currentVisualLength))
					prompt()
					os.Stdout.WriteString(recalledCommand)
					line = recalledCommand
				} else {
					if historyNavigationIndex > executor.GetHistoryLength() {
						historyNavigationIndex--
					}
					bell()
				}
			}

		case 0x7f: // BACKSPACE
			if len(line) != 0 {
				line = line[:len(line)-1]
				os.Stdout.Write([]byte{'\b', ' ', '\b'})
			}
		default:
			os.Stdout.Write(buffer)
			line += string(character)
			historyNavigationIndex = executor.GetHistoryLength()
		}
	}
}
