package main

import "fmt"

func main() {
	start := []int{18, 11, 9, 0, 5, 1}
	// start := []int{2, 1, 3}

	m := make(map[int]int)
	for i, v := range start[:len(start)-1] {
		m[v] = i + 1
	}

	last := start[len(start)-1]
	for i := len(start); i < 30000000; i++ {
		// fmt.Printf("Last: %d\n", last)

		index := m[last]
		m[last] = i

		if index == 0 {
			last = 0
		} else {
			last = i - index
		}
	}

	// fmt.Printf("%v\n", numbers)
	fmt.Printf("Letzte Nummer: %d\n", last)
}
