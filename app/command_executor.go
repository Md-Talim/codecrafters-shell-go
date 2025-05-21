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
	command1Def := parsedCommands[0]
	command2Def := parsedCommands[1]

	if len(command1Def) == 0 || len(command2Def) == 0 {
		fmt.Fprintln(finalShellIO.ErrorFile(), "shell: error, empty command in pipeline")
		return
	}

	command1Name, command1Args := command1Def[0], command1Def[1:]
	command2Name, command2Args := command2Def[0], command2Def[1:]

	// Create the pipe
	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(finalShellIO.ErrorFile(), "shell: error creating pipe: ", err)
		return
	}

	var cmd1Exec *exec.Cmd
	var cmd2Exec *exec.Cmd

	// Execute the first command (producer)
	if builtinCommand, ok := builtinCommands[command1Name]; ok {
		builtinIO := shellio.NewShellIO(pipeWriter, finalShellIO.ErrorFile())
		builtinCommand(command1Args, builtinIO)
		pipeWriter.Close()
	} else {
		cmd1Exec = exec.Command(command1Name, command1Args...)
		cmd1Exec.Stdout = pipeWriter
		cmd1Exec.Stderr = finalShellIO.ErrorFile()

		if err := cmd1Exec.Start(); err != nil {
			fmt.Fprintf(finalShellIO.ErrorFile(), "shell: error starting command %s: %v\n", command1Name, err)
			pipeWriter.Close()
			pipeReader.Close()
			return
		}
		pipeWriter.Close() // Close the write end of the pipe in the parent. cmd1 still holds it open.
	}

	// Execute the second command (consumer)
	if builtinCommand, ok := builtinCommands[command2Name]; ok {
		builtinIO := shellio.NewShellIO(finalShellIO.OutputFile(), finalShellIO.ErrorFile())
		builtinCommand(command2Args, builtinIO)
		pipeReader.Close()
	} else {
		cmd2Exec = exec.Command(command2Name, command2Args...)
		cmd2Exec.Stdin = pipeReader
		cmd2Exec.Stdout = finalShellIO.OutputFile()
		cmd2Exec.Stderr = finalShellIO.ErrorFile()

		if err := cmd2Exec.Start(); err != nil {
			fmt.Fprintf(finalShellIO.ErrorFile(), "shell: error starting command %s: %v\n", command2Name, err)
			pipeReader.Close() // Close read end as well
			if cmd1Exec != nil {
				cmd1Exec.Wait()
			}
			return
		}

		// Close the read end of the pipe in the parent. cmd2 still holds it open.
		pipeReader.Close()
	}

	// Wait for both commands to finish
	if cmd1Exec != nil {
		cmd1Exec.Wait()
	}
	if cmd2Exec != nil {
		cmd2Exec.Wait()
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
