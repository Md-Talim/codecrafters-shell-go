package main

import (
	"github.com/md-talim/codecrafters-shell-go/internal/executor"
)

func main() {
	for {
		input, result := read()

		switch result {
		case ReadResultQuit:
			return
		case ReadResultEmpty:
			continue
		case ReadResultContent:
			executor.Execute(input)
		}
	}
}
