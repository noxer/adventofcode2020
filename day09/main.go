package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	ints, err := readInts("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Zahlen: %s\n", err)
		os.Exit(1)
	}

	//	fmt.Printf("%d hat nicht die Eigenschaft!\n", findInvalidNumber(ints))
	sumSlice := findSumRange(ints, 556543474)
	if len(sumSlice) == 0 {
		fmt.Println("Keine Lösung gefunden!")
		os.Exit(1)
	}

	min, max := findMinMax(sumSlice)

	fmt.Printf("Die Lösung ist %d\n", min+max)
}

func findMinMax(sl []uint64) (uint64, uint64) {
	min := uint64(math.MaxUint64)
	max := uint64(0)

	for _, n := range sl {
		if n > max {
			max = n
		}
		if n < min {
			min = n
		}
	}

	return min, max
}

func findInvalidNumber(ints []uint64) uint64 {
	for cur := 25; cur < len(ints); cur++ {
		if !checkPreamble(ints[cur], ints[cur-25:cur]) {
			return ints[cur]
		}
	}

	return 0
}

func findSumRange(slice []uint64, targetSum uint64) []uint64 {
	var p1, p2 int
	currentSum := uint64(0)

	for {
		if p2 > len(slice) {
			break
		}

		if currentSum == targetSum {
			return slice[p1:p2]
		}
		if currentSum < targetSum {
			currentSum += slice[p2]
			p2++
		} else if currentSum > targetSum {
			currentSum -= slice[p1]
			p1++
		}
	}

	return nil
}

func checkPreamble(n uint64, preamble []uint64) bool {
	for i, a := range preamble[:len(preamble)-1] {
		for _, b := range preamble[i+1:] {
			if a+b == n {
				return true
			}
		}
	}

	return false
}

func readInts(name string) ([]uint64, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var ints []uint64
	for s.Scan() {
		i, err := strconv.ParseUint(s.Text(), 10, 64)
		if err != nil {
			continue
		}

		ints = append(ints, i)
	}

	return ints, nil
}
