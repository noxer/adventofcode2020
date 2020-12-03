package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	ints, err := readInts("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Daten: %s\n", err)
		os.Exit(1)
	}

	a, b, c, ok := checkBruteSorted3(ints, 2020)
	if !ok {
		fmt.Println("Konnte keine Lösung finden!")
		os.Exit(1)
	}

	fmt.Printf("%d + %d + %d = %d\n", a, b, c, a+b+c)
	fmt.Printf("%d * %d * %d = %d\n", a, b, c, a*b*c)
}

func checkHashmap(ints []int, targetSum int) (int, int, bool) { // O(n)
	m := make(map[int]int, len(ints))

	for _, n := range ints {
		m[n]++
	}

	for _, n := range ints {
		diff := targetSum - n
		found, ok := m[diff]
		if !ok {
			continue
		}
		if diff == n && found < 2 {
			continue
		}

		return n, diff, true
	}

	return 0, 0, false
}

func checkPointer(ints []int, targetSum int) (int, int, bool) { // O(n log n)
	sort.Ints(ints) // O(n log n)

	ptr1, ptr2 := 0, len(ints)-1

	for ptr1 != ptr2 {
		sum := ints[ptr1] + ints[ptr2]
		if sum == targetSum {
			return ints[ptr1], ints[ptr2], true
		} else if sum < targetSum {
			ptr1++
		} else {
			ptr2--
		}
	}

	return 0, 0, false
}

func checkBruteSorted(ints []int, targetSum int) (int, int, bool) { // O(n²)
	sort.Ints(ints)

	for i, n := range ints[:len(ints)-1] {
		for _, m := range ints[i+1:] {
			sum := n + m
			if sum == targetSum {
				return n, m, true
			} else if sum > targetSum {
				break
			}
		}
	}

	return 0, 0, false
}

func checkBruteSorted3(ints []int, targetSum int) (int, int, int, bool) { // O(n³)
	sort.Ints(ints)

	for i, n := range ints[:len(ints)-1] {
		for j, m := range ints[i+1:] {
			tmpSum := n + m

			if tmpSum > targetSum {
				break
			}

			for _, o := range ints[i+j+2:] {
				sum := tmpSum + o

				if sum == targetSum {
					return n, m, o, true
				} else if sum > targetSum {
					break
				}
			}
		}
	}

	return 0, 0, 0, false
}

func checkBrute(ints []int, targetSum int) (int, int, bool) { // O(n²)
	for i, n := range ints[:len(ints)-1] {
		for _, m := range ints[i+1:] {
			if n+m == targetSum {
				return n, m, true
			}
		}
	}

	return 0, 0, false
}

func readInts(name string) ([]int, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var ints []int
	for s.Scan() {
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			continue
		}

		ints = append(ints, i)
	}

	return ints, nil
}
