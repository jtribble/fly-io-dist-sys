package ioutils

import (
	"bufio"
	"io"
)

type LineReader struct {
	Lines <-chan []byte
	Errs  <-chan error
	Eof   <-chan bool
}

func NewLineReader(rd *bufio.Reader) *LineReader {
	lines := make(chan []byte)
	errs := make(chan error)
	eof := make(chan bool)
	go func() {
		for {
			switch line, err := rd.ReadBytes('\n'); err {
			case nil:
				lines <- line
			case io.EOF:
				eof <- true
			default:
				errs <- err
			}
		}
	}()
	return &LineReader{
		Lines: lines,
		Errs:  errs,
		Eof:   eof,
	}
}
