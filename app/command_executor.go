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

func executePipeCommand(parsedCommands [][]string, finalShellIO shellio.IO) {
	producerDef := parsedCommands[0]
	consumerDef := parsedCommands[1]

	if len(producerDef) == 0 || len(consumerDef) == 0 {
		fmt.Fprintln(finalShellIO.ErrorFile(), "shell: error, empty command in pipeline")
		return
	}

	numCommands := len(parsedCommands)
	pipes := make([][2]*os.File, numCommands-1)
	for i := range numCommands - 1 {
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Fprintf(finalShellIO.ErrorFile(), "shell: error creating pipe %d: %v\n", i, err)
			// Clean up any pipes created so far
			for j := range i {
				pipes[j][0].Close()
				pipes[j][1].Close()
			}
			return
		}
		pipes[i][0] = r
		pipes[i][1] = w
	}

	var runningExternalCmds []*exec.Cmd

	for i, commandDef := range parsedCommands {
		if len(commandDef) == 0 {
			fmt.Fprintln(finalShellIO.ErrorFile(), "shell: error, empty command in pipeline")
			// Close all the pipe fds before returning to prevent leaks
			for _, p := range pipes {
				p[0].Close()
				p[1].Close()
			}
			// Wait for any commands already started
			for _, cmd := range runningExternalCmds {
				cmd.Wait()
			}
			return
		}

		commandName, commandArgs := commandDef[0], commandDef[1:]

		var currentStdin *os.File
		var currentStdout *os.File

		// Determine Stdin
		if i == 0 {
			currentStdin = os.Stdin
		} else {
			currentStdin = pipes[i-1][0]
		}

		// Determine Stdout
		if i == numCommands-1 {
			currentStdout = finalShellIO.OutputFile()
		} else {
			currentStdout = pipes[i][1]
		}

		if builtinCommandExecutor, isBuiltinCommand := builtinCommands[commandName]; isBuiltinCommand {
			builtinIO := shellio.NewIO(currentStdout, finalShellIO.ErrorFile())
			builtinCommandExecutor(commandArgs, builtinIO)
		} else {
			externalCommand := exec.Command(commandName, commandArgs...)
			externalCommand.Stdin = currentStdin
			externalCommand.Stdout = currentStdout
			externalCommand.Stderr = finalShellIO.ErrorFile()

			if err := externalCommand.Start(); err != nil {
				fmt.Fprintf(finalShellIO.ErrorFile(), "shell: error starting command %s: %v\n", commandName, err)
				// Close all the pipe fds before returning to prevent leaks
				for _, p := range pipes {
					p[0].Close()
					p[1].Close()
				}
				return
			}

			runningExternalCmds = append(runningExternalCmds, externalCommand)
		}
	}

	// After all commands are configured and external ones started,
	// the parent process must close all its copies of the pipe file descriptors.
	// This signals EOF to readers when writers are done.
	for _, p := range pipes {
		p[0].Close()
		p[1].Close()
	}

	for _, cmd := range runningExternalCmds {
		cmd.Wait()
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
