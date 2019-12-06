package program

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Os interface {
	Io
}

var StdOs = osc{
	Io: &retryingIo{
		attempts: 3,
		io: &bufIo{
			reader:    bufio.NewReader(os.Stdin),
			writer:    bufio.NewWriter(os.Stdout),
			inputMsg:  "Input: ",
			outputMsg: "Output: ",
		},
	},
}

type osc struct {
	Io
}

type Io interface {
	Read() int
	Write(int)
}

type bufIo struct {
	reader    *bufio.Reader
	writer    *bufio.Writer
	inputMsg  string
	outputMsg string
}

func (bio *bufIo) Read() int {
	_, err := bio.writer.WriteString(bio.inputMsg)
	if err != nil {
		panic(fmt.Errorf("Error writing message: %w", err))
	}
	err = bio.writer.Flush()
	if err != nil {
		panic(fmt.Errorf("Error flushing message: %w", err))
	}

	in, err := bio.reader.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("Error reading message: %w", err))
	}

	in = strings.TrimSpace(in)

	i, err := strconv.Atoi(in)
	if err != nil {
		panic(fmt.Errorf("Input must be a single number: %w", err))
	}

	return i
}

func (bio *bufIo) Write(i int) {
	_, err := bio.writer.WriteString(fmt.Sprintf("%s%d\n", bio.outputMsg, i))
	if err != nil {
		panic(fmt.Errorf("Error writing output: %w", err))
	}

	err = bio.writer.Flush()
	if err != nil {
		panic(fmt.Errorf("Error flushing message: %w", err))
	}
}

type retryingIo struct {
	attempts int
	io       Io
}

var noRead = fmt.Errorf("No read")
var noWrite = fmt.Errorf("No write")

func (rio *retryingIo) tryRead() (i int, err error) {
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("Recovered panic: %v", e)
			}
		}
	}()

	return rio.io.Read(), nil
}

func (rio *retryingIo) Read() (i int) {
	err := noRead
	for attempt := 0; err != nil && attempt < rio.attempts; attempt++ {
		i, err = rio.tryRead()
		if err != nil {
			fmt.Println(err)
		}
	}

	if err != nil {
		panic(err)
	}

	return
}

func (rio *retryingIo) tryWrite(i int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("Recovered panic: %v", e)
			}
		}
	}()

	rio.io.Write(i)
	return nil
}

func (rio *retryingIo) Write(i int) {
	err := noWrite
	for attempt := 0; err != nil && attempt < rio.attempts; attempt++ {
		err = rio.tryWrite(i)
		if err != nil {
			fmt.Println(err)
		}
	}

	if err != nil {
		panic(err)
	}

	return
}
