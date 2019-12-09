package program

import (
	"fmt"
)

type Program struct {
	intcode      *Intcode
	index        int
	instructions map[int]Op
	os           Os
}

type ProgramOpt func(p *Program)

func WithIntcode(intcode *Intcode) ProgramOpt {
	return func(p *Program) {
		p.intcode = intcode
	}
}

func WithInstructions(instructions map[int]Op) ProgramOpt {
	return func(p *Program) {
		p.instructions = instructions
	}
}

func WithOs(os Os) ProgramOpt {
	return func(p *Program) {
		p.os = os
	}
}

func New(opts ...ProgramOpt) *Program {
	p := &Program{
		intcode:      NewIntcode(99),
		instructions: StdInstructions,
		os:           NewOs(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Program) Get(index int) int64 {
	return p.intcode.GetImmediate(index)
}

func (p *Program) Run() (err error) {
	for err == nil {
		err = p.Step()
	}

	if _, ok := err.(Halted); ok {
		return nil
	}

	return
}

func (p *Program) Step() (err error) {
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("Recovered panic: %v", e)
			}
		}
	}()

	if p.index < 0 {
		return Halted{}
	}

	op := p.intcode.GetImmediate(p.index)
	code := opCode(op)
	opf, ok := p.instructions[code]
	if !ok {
		return InvalidOpCode{
			opCode: code,
			index:  p.index,
		}
	}

	p.index = opf(p.intcode, p.index, op, p.os)
	return nil
}

func (p *Program) String() string {
	return p.intcode.StringIndexed(p.index)
}
