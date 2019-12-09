package program

type Op func(intcode *Intcode, index int, op int64, os Os) int

var StdInstructions = map[int]Op{
	1:  Add,
	2:  Mul,
	3:  Input,
	4:  Output,
	5:  JumpIfTrue,
	6:  JumpIfFalse,
	7:  LessThan,
	8:  Equals,
	9:  AdjustBase,
	99: Halt,
}

func Add(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 3)

	a := intcode.Get(mode[0], index+1)
	b := intcode.Get(mode[1], index+2)
	intcode.Set(mode[2], index+3, a+b)

	return index + 4
}

func Mul(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 3)

	a := intcode.Get(mode[0], index+1)
	b := intcode.Get(mode[1], index+2)
	intcode.Set(mode[2], index+3, a*b)

	return index + 4
}

func Input(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 1)

	a := os.Read()
	intcode.Set(mode[0], index+1, a)

	return index + 2
}

func Output(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 1)

	a := intcode.Get(mode[0], index+1)
	os.Write(a)

	return index + 2
}

func JumpIfTrue(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 2)

	a := intcode.Get(mode[0], index+1)
	b := intcode.Get(mode[1], index+2)

	if a != 0 {
		return int(b)
	}

	return index + 3
}

func JumpIfFalse(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 2)

	a := intcode.Get(mode[0], index+1)
	b := intcode.Get(mode[1], index+2)

	if a == 0 {
		return int(b)
	}

	return index + 3
}

func LessThan(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 3)

	a := intcode.Get(mode[0], index+1)
	b := intcode.Get(mode[1], index+2)

	if a < b {
		intcode.Set(mode[2], index+3, 1)
	} else {
		intcode.Set(mode[2], index+3, 0)
	}

	return index + 4
}

func Equals(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 3)

	a := intcode.Get(mode[0], index+1)
	b := intcode.Get(mode[1], index+2)

	if a == b {
		intcode.Set(mode[2], index+3, 1)
	} else {
		intcode.Set(mode[2], index+3, 0)
	}

	return index + 4
}

func AdjustBase(intcode *Intcode, index int, op int64, os Os) int {
	mode := modes(op, 1)

	a := intcode.Get(mode[0], index+1)
	intcode.AdjustBase(int(a))

	return index + 2
}

func Halt(intcode *Intcode, index int, op int64, os Os) int {
	return -1
}
