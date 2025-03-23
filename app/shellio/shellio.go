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
	ErrorFile() *os.File
	Close()
}

type FileRedirect struct {
	Output *os.File
	Error  *os.File
}

func (io *FileRedirect) OutputFile() *os.File {
	if io.Output != nil {
		return io.Output
	} else {
		return os.Stdout
	}
}

func (io *FileRedirect) ErrorFile() *os.File {
	if io.Error != nil {
		return io.Error
	} else {
		return os.Stderr
	}
}

func (io *FileRedirect) Close() {
	if io.Output != nil {
		io.Output.Close()
		io.Output = nil
	}
	if io.Error != nil {
		io.Error.Close()
		io.Error = nil
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

	if redirect.Descriptor == 1 {
		return &FileRedirect{Output: file, Error: nil}, true
	} else {
		return &FileRedirect{Output: nil, Error: file}, true
	}
}
