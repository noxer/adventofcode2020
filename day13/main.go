package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/deanveloper/modmath"
)

func main() {
	start, busses, err := readTimetable("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen des Fahrplans: %s\n", err)
		os.Exit(1)
	}

	minWait := 99999999999999
	minBus := 0
	for _, bus := range busses {
		if bus == -1 {
			continue
		}
		wait := findAbove(start, bus)
		if wait < minWait {
			minWait = wait
			minBus = bus
		}
	}

	fmt.Printf("Lösungsbus: %d mit %d Wartezeit\n", minBus, minWait)
	fmt.Printf("Lösung: %d\n", minBus*minWait)

	var entries []modmath.CrtEntry
	var ints []int
	for i, bus := range busses {
		if bus < 0 {
			continue
		}

		entry := modmath.CrtEntry{
			A: bus - i,
			M: bus,
		}

		entries = append(entries, entry)
		ints = append(ints, bus)
	}

	fmt.Printf("%#v\n", entries)

	solution := modmath.SolveCrtMany(entries)

	fmt.Printf("Lösung 2: %d\n", solution%lcm(ints[0], ints[1], ints[2:]...))
}

func filterTimes(times, filtered chan int, bus, offset int) {
	multiple := 0
	for time := range times {
		for time > multiple+offset {
			multiple += bus
		}

		if time == multiple+offset {
			if bus > 100 {
				fmt.Printf("ID %d: %d = %d+%d\n", bus, time, multiple, offset)
			}
			filtered <- time
		}
	}
}

func checkTime(time int, busses []int) bool {
	for i, bus := range busses {
		if findAbove(time, bus) != i {
			return false
		}
	}

	return true
}

func findAbove(start, bus int) int {
	wait := start / bus
	if start%bus > 0 {
		wait++
	}

	return wait*bus - start
}

func readTimetable(name string) (int, []int, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Scan()

	start, err := strconv.Atoi(s.Text())
	if err != nil {
		return 0, nil, err
	}

	s.Scan()
	ids := strings.Split(s.Text(), ",")
	busses := make([]int, len(ids))
	for i, id := range ids {
		if id == "x" {
			busses[i] = -1
			continue
		}

		intID, err := strconv.Atoi(id)
		if err != nil {
			return 0, nil, err
		}
		busses[i] = intID
	}

	return start, busses, nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int, ints ...int) int {
	result := a * b / gcd(a, b)
	for _, n := range ints {
		result = result * n / gcd(result, n)
	}
	return result
}
