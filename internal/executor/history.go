package executor

import (
	"bufio"
	"fmt"
	"os"

	"github.com/md-talim/codecrafters-shell-go/internal/shellio"
)

type CommandHistory struct {
	commandList []string
}

func (h *CommandHistory) length() int {
	return len(h.commandList)
}

func (h *CommandHistory) at(index int) string {
	return h.commandList[index]
}

func (h *CommandHistory) add(command string) {
	h.commandList = append(h.commandList, command)
}

func (h *CommandHistory) printAll(io shellio.IO) {
	h.printLast(len(h.commandList), io)
}

func (h *CommandHistory) printLast(limit int, io shellio.IO) {
	startIndex := len(h.commandList) - limit
	for i := startIndex; i < len(h.commandList); i++ {
		fmt.Fprintf(io.OutputFile(), "    %d  %s\n", i+1, h.commandList[i])
	}
}

func (h *CommandHistory) appendFromFile(historyFile string) {
	file, err := os.Open(historyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening history file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		h.add(line)
	}
}

func (h *CommandHistory) saveToFile(historyFile string) {
	historyBytes := []byte{}
	for _, command := range h.commandList {
		historyBytes = append(historyBytes, []byte(command+"\n")...)
	}
	if err := os.WriteFile(historyFile, historyBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing history file: %v", err)
		return
	}
}

func (h *CommandHistory) appendToFile(historyFile string) {
	historyString := ""
	lastAppendIndex := h.lastAppendIndex()
	for i := lastAppendIndex + 1; i < len(h.commandList); i++ {
		command := h.commandList[i]
		historyString += command + "\n"
	}

	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	file, err := os.OpenFile(historyFile, flags, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(historyString); err != nil {
		fmt.Fprintf(os.Stderr, "error writing history file: %v", err)
	}
}

func (h *CommandHistory) lastAppendIndex() int {
	if len(h.commandList) <= 2 {
		return -1
	}

	historyLength := len(h.commandList)
	lastAppendCommand := h.commandList[historyLength-1]
	for i := historyLength - 2; i >= 0; i-- {
		if lastAppendCommand == h.commandList[i] {
			return i
		}
	}
	return -1
}
