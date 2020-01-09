package core

import (
	"errors"
	"io"
	"os"
)

type Logger struct {
	verbose bool
	writer  io.Writer
}

//NewLogger Logs to std out if verbose is set to true
func NewLogger(verbose bool) *Logger {
	return &Logger{verbose: verbose, writer: os.Stdout}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	if l.verbose {
		return l.writer.Write(p)
	}

	return 0, errors.New("Should not be writting")
}
