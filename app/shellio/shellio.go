package shellio

import (
	"fmt"
	"os"
)

// RedirectionConfig holds the configuration for file redirection.
type RedirectionConfig struct {
	File            string
	Descriptor      int
	IsAppendEnabled bool
}

type IO interface {
	OutputFile() *os.File
	ErrorFile() *os.File
	Close()
}

func NewIO(outputFile, errorFile *os.File) IO {
	return &FileRedirect{
		outputFile: outputFile,
		errorFile:  errorFile,
	}
}

type FileRedirect struct {
	outputFile *os.File
	errorFile  *os.File
}

// OutputFile returns the output file. If it is not set, it returns os.Stdout.
func (io *FileRedirect) OutputFile() *os.File {
	if io.outputFile != nil {
		return io.outputFile
	} else {
		return os.Stdout
	}
}

// ErrorFile returns the error file. If it is not set, it returns os.Stderr.
func (io *FileRedirect) ErrorFile() *os.File {
	if io.errorFile != nil {
		return io.errorFile
	} else {
		return os.Stderr
	}
}

// Close closes any files that were opened for redirection.
func (io *FileRedirect) Close() {
	if io.outputFile != nil {
		io.outputFile.Close()
		io.outputFile = nil
	}
	if io.errorFile != nil {
		io.errorFile.Close()
		io.errorFile = nil
	}
}

func OpenIo(redirect RedirectionConfig) (IO, bool) {
	if redirect.File == "" {
		return &FileRedirect{}, false
	}

	flag := os.O_CREATE | os.O_WRONLY
	if redirect.IsAppendEnabled {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	file, err := os.OpenFile(redirect.File, flag, 0664)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shell: %s: %s\n", redirect.File, err.Error())
		return &FileRedirect{}, false
	}

	if redirect.Descriptor == 1 {
		return &FileRedirect{outputFile: file, errorFile: nil}, true
	} else {
		return &FileRedirect{outputFile: nil, errorFile: file}, true
	}
}
