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

type NegativeIndex int

func (i NegativeIndex) Error() string {
	return fmt.Sprint("Negative index (%d)", int(i))
}

type InvalidOpCode struct {
	opCode int
	index  int
}

func (i InvalidOpCode) Error() string {
	return fmt.Sprintf("Invalid op code %d at index %d", i.opCode, i.index)
}

type InvalidMode struct {
	mode  Mode
	index int
}

func (i InvalidMode) Error() string {
	return fmt.Sprintf("Invalid mode %v at index %d", i.mode, i.index)
}

type Halted struct{}

func (i Halted) Error() string {
	return "Halted"
}

var NotImplemented = fmt.Errorf("Not implemented")
