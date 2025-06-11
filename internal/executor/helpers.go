package executor

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/internal/shellio"
)

// initializePipes creates a specified number of pipes for inter-process communication.
// It returns a slice of pipes, where each pipe is represented as a [2]*os.File array.
func initializePipes(numPipes int) ([][2]*os.File, error) {
	pipes := make([][2]*os.File, numPipes)
	for i := range numPipes {
		r, w, err := os.Pipe()
		if err != nil {
			for j := range i {
				pipes[j][0].Close()
				pipes[j][1].Close()
			}
			return nil, fmt.Errorf("shell: error creating pipe %d: %v", i, err)
		}
		pipes[i] = [2]*os.File{r, w}
	}
	return pipes, nil
}

func addCommandToHistory(command string) {
	shellHistory = append(shellHistory, command)
}

func getHistoryLimit(args *[]string) (int, error) {
	var limit int
	var err error
	if len(*args) == 1 {
		stringLimit := (*args)[0]
		limit, err = strconv.Atoi(stringLimit)
		if err != nil {
			return 0, fmt.Errorf("invalid argument: %s", stringLimit)
		}
	} else {
		limit = len(shellHistory)
	}
	return limit, nil
}

func printHistory(limit int, io shellio.IO) {
	startIndex := len(shellHistory) - limit
	for i := startIndex; i < len(shellHistory); i++ {
		fmt.Fprintf(io.OutputFile(), "    %d  %s\n", i+1, shellHistory[i])
	}
}

func appendHistoryFromFile(historyFile string) {
	file, err := os.Open(historyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening history file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		addCommandToHistory(line)
	}
}

func writeHistoryToFile(historyFile string) {
	historyBytes := []byte{}
	for _, command := range shellHistory {
		historyBytes = append(historyBytes, []byte(command+"\n")...)
	}
	if err := os.WriteFile(historyFile, historyBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing history file: %v", err)
		return
	}
}

func appendHistoryToFile(historyFile string) {
	historyString := ""
	lastAppendIndex := getLastAppendIndex()
	for i := lastAppendIndex + 1; i < len(shellHistory); i++ {
		command := shellHistory[i]
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

func getLastAppendIndex() int {
	if len(shellHistory) <= 2 {
		return -1
	}

	historyLength := GetHistoryLength()
	lastAppendCommand := shellHistory[historyLength-1]
	for i := historyLength - 2; i >= 0; i-- {
		if lastAppendCommand == shellHistory[i] {
			return i
		}
	}
	return -1
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
