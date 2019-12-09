package program

import (
	"fmt"
	"strconv"
	"strings"
)

type Intcode struct {
	code []int64
	base int
}

func NewIntcode(intcode ...int64) *Intcode {
	i := &Intcode{
		code: intcode,
	}
	return i
}

func (c *Intcode) GetImmediate(index int) int64 {
	if index < 0 {
		panic(NegativeIndex(index))
	}

	if index >= len(c.code) {
		return 0
	}

	return c.code[index]
}

func (c *Intcode) GetIndirect(index int) int64 {
	return c.GetImmediate(int(c.GetImmediate(index)))
}

func (c *Intcode) GetRelative(index int) int64 {
	return c.GetImmediate(int(c.GetImmediate(index)) + c.base)
}

func (c *Intcode) Get(mode Mode, index int) int64 {
	switch mode {
	case Indirect:
		return c.GetIndirect(index)
	case Immediate:
		return c.GetImmediate(index)
	case Relative:
		return c.GetRelative(index)
	default:
		panic(InvalidMode{mode, index})
	}
}

func (c *Intcode) SetImmediate(index int, i int64) {
	if index < 0 {
		panic(NegativeIndex(index))
	}

	if index >= len(c.code) {
		c.code = append(c.code, make([]int64, index-len(c.code)+1)...)
	}

	c.code[index] = i
}

func (c *Intcode) SetIndirect(index int, i int64) {
	c.SetImmediate(int(c.GetImmediate(index)), i)
}

func (c *Intcode) SetRelative(index int, i int64) {
	c.SetImmediate(int(c.GetImmediate(index))+c.base, i)
}

func (c *Intcode) Set(mode Mode, index int, i int64) {
	switch mode {
	case Indirect:
		c.SetIndirect(index, i)
	case Immediate:
		panic(InvalidMode{mode, index})
		// c.SetImmediate(index, i)
	case Relative:
		c.SetRelative(index, i)
	default:
		panic(InvalidMode{mode, index})
	}
}

func (c *Intcode) SetBase(base int) {
	c.base = base
}

func (c *Intcode) AdjustBase(base int) {
	c.base += base
}

func (c *Intcode) String() string {
	return c.StringIndexed(-1)
}

func pad(i ...int64) int {
	max := 0
	for _, ii := range i {
		if p := len(strconv.FormatInt(ii, 10)); p > max {
			max = p
		}
	}
	return max
}

func (c *Intcode) StringIndexed(index int) string {
	cols := 10
	rows := (len(c.code) + cols - 1) / cols

	lPad := pad(int64(len(c.code)) - 1)
	cPad := pad(c.code...)
	if hPad := pad(int64(cols) - 1); hPad > cPad {
		cPad = hPad
	}

	b := &strings.Builder{}

	fmt.Fprintf(b, strings.Repeat(" ", lPad))
	for i := 0; i < cols; i++ {
		fmt.Fprintf(b, "  %0[1]*[2]d ", cPad, i)
	}

	for row := 0; row < rows; row++ {
		fmt.Fprintf(b, "\n%0[1]*[2]d", lPad, row*10)
		for i := 0; i < cols; i++ {
			cIndex := row*10 + i
			if cIndex >= len(c.code) {
				break
			}

			if cIndex == index {
				fmt.Fprintf(b, " [%[1]*[2]d]", cPad, c.code[cIndex])
			} else {
				fmt.Fprintf(b, "  %[1]*[2]d ", cPad, c.code[cIndex])
			}
		}
	}

	return b.String()
}
