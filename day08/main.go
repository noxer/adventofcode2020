package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	program, err := readCode("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen des Programms: %s\n", err)
		os.Exit(1)
	}

	start := time.Now()
	acc := executeTryPatch(program)
	fmt.Printf("Das hat %s gedauert\n", time.Since(start))

	fmt.Printf("Acc ist %d\n", acc)
}

func executeTryPatch(program []Instruction) int {
	for i, inst := range program {
		switch v := inst.(type) {
		case Acc:
			continue
		case Nop:
			program[i] = Jmp(v)
		case Jmp:
			program[i] = Nop(v)
		}

		acc, ok := executeUntilLoop(program)
		if ok {
			return acc
		}

		program[i] = inst
	}

	return 0
}

func executeUntilLoop(program []Instruction) (int, bool) {
	state := &State{Instructions: program}

	visited := make([]bool, len(program))
	for !visited[state.NextInstruction] {
		visited[state.NextInstruction] = true

		state.Instructions[state.NextInstruction].Execute(state)
		if state.NextInstruction >= len(state.Instructions) {
			return state.Accumulator, state.NextInstruction == len(state.Instructions)
		}
	}

	return state.Accumulator, false
}

// State ...
type State struct {
	Accumulator     int
	Instructions    []Instruction
	NextInstruction int
}

// Instruction ...
type Instruction interface {
	Execute(state *State)
}

// Nop ...
type Nop int

// Execute ...
func (Nop) Execute(state *State) {
	state.NextInstruction++
}

// Acc ...
type Acc int

// Execute ...
func (a Acc) Execute(state *State) {
	state.Accumulator += int(a)
	state.NextInstruction++
}

// Jmp ...
type Jmp int

// Execute ...
func (j Jmp) Execute(state *State) {
	state.NextInstruction += int(j)
}

func readCode(name string) ([]Instruction, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var program []Instruction
	for s.Scan() {
		line := s.Text()

		if len(line) < 4 {
			continue
		}

		switch line[:3] {

		case "nop":
			i, err := strconv.Atoi(line[4:])
			if err != nil {
				return nil, err
			}
			program = append(program, Nop(i))

		case "acc":
			i, err := strconv.Atoi(line[4:])
			if err != nil {
				return nil, err
			}
			program = append(program, Acc(i))

		case "jmp":
			i, err := strconv.Atoi(line[4:])
			if err != nil {
				return nil, err
			}
			program = append(program, Jmp(i))

		default:
			return nil, fmt.Errorf("unknown instruction: %s", line)

		}
	}

	return program, nil
}
