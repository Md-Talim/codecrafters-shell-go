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

func executePipeCommand(parsedCommands [][]string, shellio shellio.IO) {
	command1Def := parsedCommands[0]
	command2Def := parsedCommands[1]

	if len(command1Def) == 0 || len(command2Def) == 0 {
		fmt.Fprintln(shellio.ErrorFile(), "shell: error, empty command in pipeline")
		return
	}

	command1Name, command1Args := command1Def[0], command1Def[1:]
	command2Name, command2Args := command2Def[0], command2Def[1:]

	// Create the pipe
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(shellio.ErrorFile(), "shell: error creating pipe: ", err)
		return
	}

	// First command
	cmd1 := exec.Command(command1Name, command1Args...)
	cmd1.Stdout = w
	cmd1.Stderr = shellio.ErrorFile()

	// Second command
	cmd2 := exec.Command(command2Name, command2Args...)
	cmd2.Stdin = r
	cmd2.Stdout = shellio.OutputFile()
	cmd2.Stderr = shellio.ErrorFile()

	if err := cmd1.Start(); err != nil {
		fmt.Fprintf(shellio.ErrorFile(), "shell: error starting command %s: %v\n", command1Name, err)
		w.Close()
		r.Close()
		return
	}
	// Close the write end of the pipe in the parent. cmd1 still holds it open.
	// This is crucial so that cmd2 receives EOF when cmd1 finishes.
	w.Close()

	if err := cmd2.Start(); err != nil {
		fmt.Fprintf(shellio.ErrorFile(), "shell: error starting command %s: %v\n", command2Name, err)
		r.Close()   // Close read end as well
		cmd1.Wait() // Attempt to reap cmd1 if it started
		return
	}

	// Close the read end of the pipe in the parent. cmd2 still holds it open.
	r.Close()

	// Wait for both commands to finish
	if err := cmd1.Wait(); err != nil {
		fmt.Fprintf(shellio.ErrorFile(), "shell: error waiting for command %s: %v\n", command1Name, err)
	}
	if err := cmd2.Wait(); err != nil {
		fmt.Fprintf(shellio.ErrorFile(), "shell: error waiting for command %s: %v\n", command2Name, err)
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
