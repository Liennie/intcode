package program

import (
	"fmt"
)

var (
	_ error = NegativeIndex(0)
	_ error = InvalidOpCode{}
	_ error = InvalidMode{}
	_ error = Halted{}
	_ error = NotImplemented
)

type NegativeIndex int64

func (i NegativeIndex) Error() string {
	return fmt.Sprint("Negative index (%d)", int64(i))
}

type InvalidOpCode struct {
	opCode int64
	index  int64
}

func (i InvalidOpCode) Error() string {
	return fmt.Sprintf("Invalid op code %d at index %d", i.opCode, i.index)
}

type InvalidMode struct {
	mode  Mode
	index int64
}

func (i InvalidMode) Error() string {
	return fmt.Sprintf("Invalid mode %v at index %d", i.mode, i.index)
}

type Halted struct{}

func (i Halted) Error() string {
	return "Halted"
}

var NotImplemented = fmt.Errorf("Not implemented")

var noRead = fmt.Errorf("No read")
var noWrite = fmt.Errorf("No write")
var noData = fmt.Errorf("No data")
