package util

import "sort"

func copyInts(ints []int) []int {
	newInts := make([]int, len(ints))
	copy(newInts, ints)
	return newInts
}

func reverseInts(ints []int) []int {
	ints = copyInts(ints)

	i := 0
	j := len(ints) - 1
	for i < j {
		ints[i], ints[j] = ints[j], ints[i]

		i++
		j--
	}
	return ints
}

func PermInts(ints ...int) <-chan []int {
	ints = copyInts(ints)
	n := len(ints)

	sort.Ints(ints)
	ch := make(chan []int)

	go func() {
		defer close(ch)

		ch <- copyInts(ints)

		for {
			k := n - 2
			for k >= 0 {
				if ints[k] < ints[k+1] {
					break
				}

				k--
			}
			if k < 0 {
				return
			}

			l := n - 1
			for l > k {
				if ints[l] > ints[k] {
					break
				}

				l--
			}

			ints[k], ints[l] = ints[l], ints[k]

			ints = append(ints[:k+1], reverseInts(ints[k+1:])...)

			ch <- copyInts(ints)
		}
	}()

	return ch
}
