package program

type Mode int

const (
	Indirect  Mode = 0
	Immediate Mode = 1
	Relative  Mode = 2
)

func modes(op int64, count int) []Mode {
	op /= 100

	res := make([]Mode, count)

	for i := 0; i < count; i++ {
		res[i] = Mode(op % 10)
		op /= 10
	}

	return res
}

func opCode(op int64) int {
	return int(op % 100)
}
