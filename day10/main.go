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
		fmt.Printf("Fehler beim Einlesen der Adaper Liste: %s\n", err)
		os.Exit(1)
	}

	ints = append(ints, 0)
	sort.Ints(ints)
	ints = append(ints, ints[len(ints)-1]+3)

	//	step1, step3 := findSteps(ints)
	//	fmt.Printf("Lösung: %d\n", step1*step3)

	runs := splitInts(ints)
	fmt.Printf("%v\n", runs)

	prod := uint64(1)
	for _, run := range runs {
		prod *= calcNaive(run[0], run[1:])
	}

	fmt.Printf("Anzahl der Möglichkeiten: %d\n", prod)
}

func findSteps(ints []int) (int, int) {
	step1 := 0
	step3 := 0
	for i, a := range ints[:len(ints)-1] {
		diff := ints[i+1] - a
		switch diff {
		case 1:
			step1++
		case 3:
			step3++
		}
	}

	return step1, step3
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

func splitInts(ints []int) [][]int {
	var result [][]int

	start := 0
	for i, a := range ints[:len(ints)-1] {
		diff := ints[i+1] - a
		if diff != 3 {
			continue
		}

		result = append(result, ints[start:i+1])
		start = i + 1
	}

	return result
}

func calcNaive(first int, remaining []int) uint64 {
	if len(remaining) == 0 {
		return 1
	}

	sum := uint64(0)
	for remaining[0]-first <= 3 {
		sum += calcNaive(remaining[0], remaining[1:])
		remaining = remaining[1:]
		if len(remaining) == 0 {
			break
		}
	}

	return sum
}
