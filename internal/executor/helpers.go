package executor

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
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

func loadHistoryFromHISTFILE() {
	if historyFileName, isPresent := os.LookupEnv("HISTFILE"); isPresent {
		history.appendFromFile(historyFileName)
	}
}

func writeHistoryToHISTFILE() {
	if historyFileName, isPresent := os.LookupEnv("HISTFILE"); isPresent {
		history.saveToFile(historyFileName)
	}
}

func parseHistoryLimit(stringLimit string) (int, error) {
	var limit int
	var err error
	limit, err = strconv.Atoi(stringLimit)
	if err != nil {
		return 0, fmt.Errorf("invalid argument: %s", stringLimit)
	}
	return limit, nil
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
