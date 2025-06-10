package executor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/md-talim/codecrafters-shell-go/internal/shellio"
)

func (pr *PipelineRunner) determineStageIO(commandIndex, numTotalCommands int) (stdin, stdout *os.File) {
	// Determine Stdin
	if commandIndex == 0 {
		stdin = os.Stdin
	} else {
		stdin = pr.pipes[commandIndex-1][0]
	}

	// Determine Stdout
	if commandIndex == numTotalCommands-1 {
		stdout = pr.finalShellIO.OutputFile()
	} else {
		stdout = pr.pipes[commandIndex][1]
	}
	return stdin, stdout
}

// cleanupPipelineResources closes all pipes and waits for all running external commands to finish.
func (pr *PipelineRunner) cleanupPipelineResources() {
	for _, p := range pr.pipes {
		if p[0] != nil {
			p[0].Close()
		}
		if p[1] != nil {
			p[1].Close()
		}
	}

	for _, cmd := range pr.runningExternalCommands {
		cmd.Wait()
	}
}

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
