package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	commands, err := readCommands("input.txt")
	if err != nil {
		fmt.Printf("Konnte die Liste von Befehlen nicht einlesen: %s\n", err)
		os.Exit(1)
	}

	s := state{
		mem: make(map[uint64]uint64),
	}

	for _, cmd := range commands {
		cmd.Execute2(&s)
	}

	sum := uint64(0)
	for _, value := range s.mem {
		sum += value
	}

	fmt.Printf("Ergebnis: %d\n", sum)
}

type state struct {
	mask0 uint64
	mask1 uint64
	maskX []uint64
	mem   map[uint64]uint64
}

type command interface {
	Execute(*state)
	Execute2(*state)
}

type mask string

func (m mask) Execute2(s *state) {
	var mask0, mask1 uint64

	maskX := make([]uint64, 1, 1<<strings.Count(string(m), "X"))

	for _, c := range m {
		mask0 <<= 1
		mask1 <<= 1
		shiftMasks(maskX)

		switch c {
		case '0':
			mask0 |= 1

		case '1':
			mask0 |= 1
			mask1 |= 1

		case 'X':
			var added []uint64
			for _, x := range maskX {
				added = append(added, x|1)
			}
			maskX = append(maskX, added...)
		}
	}

	s.mask0 = mask0
	s.mask1 = mask1
	s.maskX = maskX

	fmt.Printf("Mask:  %s\n", m)
	fmt.Printf("Mask0: %36b\n", mask0)
	fmt.Printf("Mask1: %36b\n", mask1)
	for i, x := range maskX {
		fmt.Printf("MaskX: %36b (%d)\n", x, i)
	}
}

func shiftMasks(masks []uint64) {
	for i, m := range masks {
		masks[i] = m << 1
	}
}

func (m mask) Execute(s *state) {
	var mask0, mask1 uint64
	for _, c := range m {
		mask0 <<= 1
		mask1 <<= 1

		switch c {
		case 'X':
			mask0 |= 1

		case '1':
			mask0 |= 1
			mask1 |= 1
		}
	}

	s.mask0 = mask0
	s.mask1 = mask1

	//	fmt.Printf("Mask:  %s\n", m)
	//	fmt.Printf("Mask0: %36b\n", mask0)
	//	fmt.Printf("Mask1: %36b\n", mask1)
}

type mem struct {
	addr  uint64
	value uint64
}

func (m mem) Execute(s *state) {
	val := m.value & s.mask0
	val = val | s.mask1

	s.mem[m.addr] = val
}

func (m mem) Execute2(s *state) {
	addr := m.addr & s.mask0
	addr = addr | s.mask1

	for _, x := range s.maskX {
		s.mem[addr|x] = m.value
	}
}

var (
	matchMask = regexp.MustCompile(`^mask = ([01X]+)$`)
	matchMem  = regexp.MustCompile(`^mem\[(\d+)\] = (\d+)$`)
)

func readCommands(name string) ([]command, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var commands []command
	for s.Scan() {
		line := s.Text()

		if matchMask.MatchString(line) {
			matches := matchMask.FindStringSubmatch(line)
			commands = append(commands, mask(matches[1]))
		} else if matchMem.MatchString(line) {
			matches := matchMem.FindStringSubmatch(line)
			m := mem{}
			m.addr, _ = strconv.ParseUint(matches[1], 10, 36)
			m.value, _ = strconv.ParseUint(matches[2], 10, 36)
			commands = append(commands, m)
		} else {
			fmt.Printf("Konnte %s nicht parsen!\n", line)
		}
	}

	return commands, nil
}
