package program

type Mode int

const (
	Indirect  Mode = 0
	Immediate Mode = 1
)

func modes(op int, count int) []Mode {
	op /= 100

	res := make([]Mode, count)

	for i := 0; i < count; i++ {
		res[i] = Mode(op % 10)
		op /= 10
	}

	return res
}

func opCode(op int) int {
	return op % 100
}
