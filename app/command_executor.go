package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/app/shellio"
)

func executeExternalCommand(command string, args []string, io shellio.IO) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = io.OutputFile()
	cmd.Stderr = io.ErrorFile()
	cmd.Run()
}

// handleProducerInPipeline handles the producer command in a pipeline.
// It executes the command and writes its output to the pipe.
// It returns the exec.Cmd if an external command was started, and an error if one occurred.
func handleProducerInPipeline(producerDef []string, pipeWriteEnd *os.File, errorFile *os.File) (*exec.Cmd, error) {
	commandName, commandArgs := producerDef[0], producerDef[1:]
	if executeBuiltinCommand, isBuiltinCommand := builtinCommands[commandName]; isBuiltinCommand {
		builtinIO := shellio.NewIO(pipeWriteEnd, errorFile)
		executeBuiltinCommand(commandArgs, builtinIO)
		pipeWriteEnd.Close()
		return nil, nil
	}

	externalCommand := exec.Command(commandName, commandArgs...)
	externalCommand.Stdout = pipeWriteEnd
	externalCommand.Stderr = errorFile
	if err := externalCommand.Start(); err != nil {
		return nil, fmt.Errorf("error starting producer %s: %v", commandName, err)
	}

	pipeWriteEnd.Close()
	return externalCommand, nil
}

// handleConsumerInPipeline handles the consumer command in a pipeline.
// It executes the command and reads its input from the pipe.
// It returns the exec.Cmd if an external command was started, and an error if one occurred.
func handleConsumerInPipeline(consumerDef []string, pipeReadEnd *os.File, io shellio.IO) (*exec.Cmd, error) {
	commandName, commandArgs := consumerDef[0], consumerDef[1:]
	if executeBuiltinCommand, isBuiltinCommand := builtinCommands[commandName]; isBuiltinCommand {
		builtinIO := shellio.NewIO(io.OutputFile(), io.ErrorFile())
		executeBuiltinCommand(commandArgs, builtinIO)
		return nil, nil
	}

	externalCommand := exec.Command(commandName, commandArgs...)
	externalCommand.Stdin = pipeReadEnd
	externalCommand.Stdout = io.OutputFile()
	externalCommand.Stderr = io.ErrorFile()
	if err := externalCommand.Start(); err != nil {
		return nil, fmt.Errorf("error starting consumer %s: %v", commandName, err)
	}

	pipeReadEnd.Close()
	return externalCommand, nil
}

func executePipeCommand(parsedCommands [][]string, finalShellIO shellio.IO) {
	producerDef := parsedCommands[0]
	consumerDef := parsedCommands[1]

	if len(producerDef) == 0 || len(consumerDef) == 0 {
		fmt.Fprintln(finalShellIO.ErrorFile(), "shell: error, empty command in pipeline")
		return
	}

	pipeReadEnd, pipeWriteEnd, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(finalShellIO.ErrorFile(), "shell: error creating pipe: ", err)
		return
	}

	var (
		producerCommand, consumerCommand *exec.Cmd
		producerErr, consumerErr         error
	)

	producerCommand, producerErr = handleProducerInPipeline(producerDef, pipeWriteEnd, finalShellIO.ErrorFile())
	if producerErr != nil {
		fmt.Fprintln(finalShellIO.ErrorFile(), producerErr.Error())
		pipeWriteEnd.Close()
		pipeReadEnd.Close()
		return
	}

	consumerCommand, consumerErr = handleConsumerInPipeline(consumerDef, pipeReadEnd, finalShellIO)
	if consumerErr != nil {
		fmt.Fprintln(finalShellIO.ErrorFile(), consumerErr.Error())
		pipeReadEnd.Close()
		if producerCommand != nil {
			producerCommand.Wait()
		}
		return
	}

	if producerCommand != nil {
		producerCommand.Wait()
	}
	if consumerCommand != nil {
		consumerCommand.Wait()
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
		executeExternalCommand(commandName, commandArgs, shellio)
	} else {
		fmt.Fprintf(shellio.OutputFile(), "%s: command not found\n", commandName)
	}
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
