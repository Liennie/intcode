package program

import (
	"fmt"
	"strconv"
	"strings"
)

type Intcode []int

func NewIntcode(intcode ...int) *Intcode {
	i := Intcode(intcode)
	return &i
}

func (c *Intcode) GetImmediate(index int) int {
	if index < 0 {
		panic(NegativeIndex(index))
	}

	if index >= len(*c) {
		return 0
	}

	return (*c)[index]
}

func (c *Intcode) GetIndirect(index int) int {
	return c.GetImmediate(c.GetImmediate(index))
}

func (c *Intcode) Get(mode Mode, index int) int {
	switch mode {
	case Indirect:
		return c.GetIndirect(index)
	case Immediate:
		return c.GetImmediate(index)
	default:
		panic(InvalidMode{mode, index})
	}
}

func (c *Intcode) SetImmediate(index, i int) {
	if index < 0 {
		panic(NegativeIndex(index))
	}

	if index >= len(*c) {
		*c = append(*c, make([]int, index-len(*c)+1)...)
	}

	(*c)[index] = i
}

func (c *Intcode) SetIndirect(index, i int) {
	c.SetImmediate(c.GetImmediate(index), i)
}

func (c *Intcode) Set(mode Mode, index, i int) {
	switch mode {
	case Indirect:
		c.SetIndirect(index, i)
	case Immediate:
		panic(InvalidMode{mode, index})
		// c.SetImmediate(index, i)
	default:
		panic(InvalidMode{mode, index})
	}
}

func (c *Intcode) String() string {
	return c.StringIndexed(-1)
}

func pad(i ...int) int {
	max := 0
	for _, ii := range i {
		if p := len(strconv.Itoa(ii)); p > max {
			max = p
		}
	}
	return max
}

func (c *Intcode) StringIndexed(index int) string {
	cols := 10
	rows := (len(*c) + cols - 1) / cols

	lPad := pad(len(*c) - 1)
	cPad := pad((*c)...)
	if hPad := pad(cols - 1); hPad > cPad {
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
			if cIndex >= len(*c) {
				break
			}

			if cIndex == index {
				fmt.Fprintf(b, " [%[1]*[2]d]", cPad, (*c)[cIndex])
			} else {
				fmt.Fprintf(b, "  %[1]*[2]d ", cPad, (*c)[cIndex])
			}
		}
	}

	return b.String()
}
