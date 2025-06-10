package executor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/md-talim/codecrafters-shell-go/internal/shellio"
)

type PipelineRunner struct {
	finalShellIO            shellio.IO
	parsedCommands          [][]string
	pipes                   [][2]*os.File
	runningExternalCommands []*exec.Cmd
}

func newPipelineRunner(parsedCommands [][]string, finalShellIO shellio.IO) *PipelineRunner {
	numCommands := len(parsedCommands)
	pipes, err := initializePipes(numCommands - 1)
	if err != nil {
		fmt.Fprintln(finalShellIO.ErrorFile(), err)
		return nil
	}
	var runningExternalCommands []*exec.Cmd
	return &PipelineRunner{
		finalShellIO:            finalShellIO,
		parsedCommands:          parsedCommands,
		pipes:                   pipes,
		runningExternalCommands: runningExternalCommands,
	}
}

func (pr *PipelineRunner) run() {
	for i, commandDef := range pr.parsedCommands {
		if len(commandDef) == 0 {
			fmt.Fprintln(pr.finalShellIO.ErrorFile(), "shell: error, empty command in pipeline")
			pr.cleanupPipelineResources()
			return
		}

		currentStdin, currentStdout := pr.determineStageIO(i, len(pr.parsedCommands))
		command, err := pr.executePipelineStage(commandDef, currentStdin, currentStdout, pr.finalShellIO.ErrorFile())
		if err != nil {
			fmt.Fprintln(pr.finalShellIO.ErrorFile(), err)
			pr.cleanupPipelineResources()
		}

		if command != nil {
			pr.runningExternalCommands = append(pr.runningExternalCommands, command)
		}
	}

	pr.cleanupPipelineResources()
}

func (pr *PipelineRunner) executePipelineStage(commandDef []string, stdin, stdout, stderr *os.File) (*exec.Cmd, error) {
	commandName, commandArgs := commandDef[0], commandDef[1:]
	if builtinCommandExecutor, isBuiltinCommand := builtinCommands[commandName]; isBuiltinCommand {
		builtinIO := shellio.NewIO(stdout, stderr)
		builtinCommandExecutor(commandArgs, builtinIO)
		return nil, nil
	}

	externalCommand := exec.Command(commandName, commandArgs...)
	externalCommand.Stdin = stdin
	externalCommand.Stdout = stdout
	externalCommand.Stderr = stderr

	if err := externalCommand.Start(); err != nil {
		return nil, fmt.Errorf("shell: error starting command %s: %v", commandName, err)
	}

	return externalCommand, nil
}

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
