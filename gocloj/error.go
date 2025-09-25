package gocloj

import (
	"fmt"
)

type Error interface {
	error
	Line() int
	File() string
}

type perror struct {
	msg  string
	line int
	file string
}

func (err perror) Error() string {
	return fmt.Sprintf("%s:%d: %s", err.file, err.line, err.msg)
}

func (err perror) Line() int {
	return err.line
}

func (err perror) File() string {
	return err.file
}

func NewError(msg string, line int, file string) Error {
	return perror{
		msg:  msg,
		line: line,
		file: file,
	}
}
