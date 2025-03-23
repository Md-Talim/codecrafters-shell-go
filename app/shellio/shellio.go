package shellio

import (
	"fmt"
	"os"
)

type RedirectionConfig struct {
	File       string
	Descriptor int
}

type IO interface {
	OutputFile() *os.File
	Close()
}

type FileRedirect struct {
	Output *os.File
}

func (io *FileRedirect) OutputFile() *os.File {
	if io.Output != nil {
		return io.Output
	} else {
		return os.Stdout
	}
}

func (io *FileRedirect) Close() {
	if io.Output != nil {
		io.Output.Close()
		io.Output = nil
	}
}

func OpenIo(redirect RedirectionConfig) (IO, bool) {
	if redirect.File == "" {
		return &FileRedirect{}, false
	}

	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	file, err := os.OpenFile(redirect.File, flag, 0664)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shell: %s: %s\n", redirect.File, err.Error())
		return &FileRedirect{}, false
	}

	return &FileRedirect{Output: file}, true
}
