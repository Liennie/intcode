package program

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Os interface {
	Io
}

type osc struct {
	Io
}

type OsOpt func(o *osc)

func WithIo(i Io) OsOpt {
	return func(o *osc) {
		o.Io = i
	}
}

func NewOs(opts ...OsOpt) Os {
	o := osc{
		Io: RetryingIo(
			3,
			BufIo(
				os.Stdin,
				os.Stdout,
				"Input: ",
				"Output: ",
			),
		),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

type Io interface {
	Read() int64
	Write(int64)
}

type nilIo struct{}

var _ Io = nilIo{}

func (nilIo) Read() int64 { panic(noData) }
func (nilIo) Write(int64) {}

type bufIo struct {
	reader    *bufio.Reader
	writer    *bufio.Writer
	inputMsg  string
	outputMsg string
}

func BufIo(reader io.Reader, writer io.Writer, inputMsg, outputMsg string) Io {
	return &bufIo{
		reader:    bufio.NewReader(reader),
		writer:    bufio.NewWriter(writer),
		inputMsg:  inputMsg,
		outputMsg: outputMsg,
	}
}

func (bio *bufIo) Read() int64 {
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

	i, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		panic(fmt.Errorf("Input must be a single number: %w", err))
	}

	return i
}

func (bio *bufIo) Write(i int64) {
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

func RetryingIo(attempts int, io Io) Io {
	if io == nil {
		io = nilIo{}
	}
	return &retryingIo{
		attempts: attempts,
		io:       io,
	}
}

func (rio *retryingIo) tryRead() (i int64, err error) {
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

func (rio *retryingIo) Read() (i int64) {
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

func (rio *retryingIo) tryWrite(i int64) (err error) {
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

func (rio *retryingIo) Write(i int64) {
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

type echoIo struct {
	buffer []int64
}

func EchoIo(initial ...int64) Io {
	return &echoIo{
		buffer: initial,
	}
}

func (eio *echoIo) Read() int64 {
	if len(eio.buffer) == 0 {
		panic(noData)
	}

	i := eio.buffer[0]
	eio.buffer = eio.buffer[1:]
	return i
}

func (eio *echoIo) Write(i int64) {
	eio.buffer = append(eio.buffer, i)
}

type chanIo struct {
	ch chan int64
}

func ChanIo(ch chan int64) Io {
	return chanIo{
		ch: ch,
	}
}

func (cio chanIo) Read() int64 {
	return <-cio.ch
}

func (cio chanIo) Write(i int64) {
	cio.ch <- i
}

type chainIo struct {
	inch <-chan int64
	ouch chan<- int64
}

func ChainIo(input <-chan int64, output chan<- int64) Io {
	return &chainIo{
		inch: input,
		ouch: output,
	}
}

func (cio *chainIo) Read() int64 {
	return <-cio.inch
}

func (cio *chainIo) Write(i int64) {
	cio.ouch <- i
}

type verboseIo struct {
	io    Io
	extra string
}

func VerboseIo(io Io, extra map[string]interface{}) Io {
	if io == nil {
		io = nilIo{}
	}

	fields := []string{}
	for k, v := range extra {
		fields = append(fields, fmt.Sprintf("%s: %v", k, v))
	}
	x := strings.Join(fields, ",")

	return &verboseIo{
		io:    io,
		extra: x,
	}
}

func (vio *verboseIo) Read() int64 {
	i := vio.io.Read()
	fmt.Printf("Read(%s): %d\n", vio.extra, i)
	return i
}

func (vio *verboseIo) Write(i int64) {
	fmt.Printf("Write(%s): %d\n", vio.extra, i)
	vio.io.Write(i)
}
