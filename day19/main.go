package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	rules, inputs, err := readInput("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Daten: %s\n", err)
		os.Exit(1)
	}

	///*
	rules[8], _ = parseRule("8: 42 | 42 8")
	rules[11], _ = parseRule("11: 42 31 | 42 11 31")
	//*/

	sum := 0
out:
	for _, input := range inputs {
		set := rules[0].Match(rules, input)

		for i := range set {
			if i == len(input) {
				sum++
				continue out
			}
		}
	}

	fmt.Printf("%d Eingaben stimmen mit Regel 0 überein...\n", sum)
}

// Rule ...
type Rule struct {
	ID       int
	Subrules [][]int
	Char     byte
}

// IntSet ...
type IntSet map[int]struct{}

// Has ...
func (s IntSet) Has(i int) bool {
	_, ok := s[i]
	return ok
}

// Put ...
func (s IntSet) Put(i int) {
	s[i] = struct{}{}
}

// Add ...
func (s IntSet) Add(i int) IntSet {
	n := IntSet{}
	for j := range s {
		n.Put(i + j)
	}
	return n
}

// Merge ...
func (s IntSet) Merge(o IntSet) {
	for i := range o {
		s[i] = struct{}{}
	}
}

// AddSet ...
func (s IntSet) AddSet(o IntSet) IntSet {
	r := IntSet{}
	for i := range s {
		for j := range o {
			r[i+j] = struct{}{}
		}
	}

	return r
}

// Match ...
func (r *Rule) Match(allRules map[int]*Rule, str string) IntSet {
	if len(str) == 0 {
		return nil
	}

	// match terminal
	if len(r.Subrules) == 0 {
		if str[0] == r.Char {
			return IntSet{1: struct{}{}}
		}
		return nil
	}

	// match subrules
	results := IntSet{}
	for _, rules := range r.Subrules {
		set := applyRules(allRules, rules, str)
		results.Merge(set)
	}

	return results
}

func applyRules(rules map[int]*Rule, ruleList []int, str string) IntSet {
	if len(ruleList) == 0 {
		return nil
	}

	sums := rules[ruleList[0]].Match(rules, str)
	if len(sums) == 0 {
		return nil
	}

	if len(ruleList) < 2 {
		return sums
	}

	result := IntSet{} // nicht IntSet{1: struct{}{}}

	count := 0
	for i := range sums {
		tmpRes := applyRules(rules, ruleList[1:], str[i:])
		if len(tmpRes) == 0 {
			continue
		}
		count++
		result.Merge(tmpRes.Add(i)) // nicht result.AddSet!
	}

	if count == 0 {
		return nil
	}

	return result
}

func parseRule(line string) (*Rule, error) {
	parts := strings.SplitN(line, ": ", 2)
	if len(parts) != 2 {
		return nil, errors.New("unexpected format")
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	r := &Rule{
		ID: id,
	}

	if parts[1][0] == '"' {
		r.Char = parts[1][1]
		return r, nil
	}

	submatches := strings.Split(parts[1], " | ")
	r.Subrules = make([][]int, len(submatches))
	for i, match := range submatches {
		ints := strings.Split(match, " ")
		r.Subrules[i] = make([]int, len(ints))

		for j, a := range ints {
			r.Subrules[i][j], err = strconv.Atoi(a)
			if err != nil {
				return nil, err
			}
		}
	}

	return r, nil
}

func readInput(name string) (map[int]*Rule, []string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	rules := make(map[int]*Rule)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}

		r, err := parseRule(line)
		if err != nil {
			fmt.Printf("Fehler beim Parsen von %s: %s\n", line, err)
			continue
		}

		rules[r.ID] = r
	}

	var inputs []string
	for s.Scan() {
		inputs = append(inputs, s.Text())
	}

	return rules, inputs, nil
}
