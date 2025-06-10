package executor

import (
	"fmt"
	"os/exec"

	"github.com/md-talim/codecrafters-shell-go/internal/parser"
	"github.com/md-talim/codecrafters-shell-go/internal/shellio"
)

func Execute(input string) {
	addCommandToHistory(input)
	p := parser.NewParser(input)
	parsedCommands, redirectionConfig := p.Parse()

	if len(parsedCommands) == 0 {
		return
	}

	finalIO, isRedirected := shellio.OpenIo(redirectionConfig)
	if isRedirected {
		defer finalIO.Close()
	}

	if len(parsedCommands) == 1 {
		executeSingleCommand(parsedCommands[0], finalIO)
	} else {
		executePipelines(parsedCommands, finalIO)
	}
}

func executeSingleCommand(command []string, finalShellIO shellio.IO) {
	if len(command) == 0 {
		return
	}
	commandName := command[0]
	commandArgs := command[1:]

	if builtinCommandExecutor, isBuiltinCommand := builtinCommands[commandName]; isBuiltinCommand {
		builtinCommandExecutor(commandArgs, finalShellIO)
	} else if _, ok := findPath(commandName); ok {
		executeExternalCommand(commandName, commandArgs, finalShellIO)
	} else {
		fmt.Fprintf(finalShellIO.OutputFile(), "%s: command not found\n", commandName)
	}
}

func executeExternalCommand(command string, args []string, io shellio.IO) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = io.OutputFile()
	cmd.Stderr = io.ErrorFile()
	cmd.Run()
}

func executePipelines(parsedCommands [][]string, finalShellIO shellio.IO) {
	pipelineRunner := newPipelineRunner(parsedCommands, finalShellIO)
	pipelineRunner.run()
}
