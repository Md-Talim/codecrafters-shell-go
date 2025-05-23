package executor

import (
	"fmt"
	"os"
	"strconv"
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

func getHistoryLimit(args []string) (int, error) {
	if len(args) > 1 {
		return 0, fmt.Errorf("history: too many arguments")
	}
	if len(shellHistory) == 0 {
		return 0, fmt.Errorf("no commands in history")
	}

	var limit int
	var err error
	if len(args) == 1 {
		limit, err = strconv.Atoi(args[0])
		if err != nil {
			return 0, fmt.Errorf("invalid argument: %s", args[0])
		}
	} else {
		limit = len(shellHistory)
	}
	return limit, nil
}
