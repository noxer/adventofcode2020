package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	rules, myTicket, tickets, err := readFile("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Datei: %s\n", err)
		os.Exit(1)
	}

	validTickets := filterTickets(tickets, rules)

	matchingRules := make([]map[Rule]struct{}, len(myTicket))
	for i := 0; i < len(myTicket); i++ {
		matchingRules[i] = identifyField(rules, validTickets, i)
	}

	for isolateSingles(matchingRules) != len(matchingRules) {
	}

	prod := 1
	for i, mr := range matchingRules {
		for r := range mr {
			if strings.HasPrefix(r.Name, "departure") {
				prod *= myTicket[i]
			}
		}
	}

	fmt.Printf("Ergebnis: %d\n", prod)
}

func isolateSingles(matchingRules []map[Rule]struct{}) int {
	var singles []Rule
	for _, mr := range matchingRules {
		if len(mr) == 1 {
			for r := range mr {
				singles = append(singles, r)
			}
		}
	}

	for _, mr := range matchingRules {
		if len(mr) == 1 {
			continue
		}

		for _, single := range singles {
			delete(mr, single)
		}
	}

	return len(singles)
}

func identifyField(rules []Rule, tickets []Ticket, index int) map[Rule]struct{} {
	m := ruleSliceToMap(rules)

	for _, ticket := range tickets {
		v := ticket[index]
		for r := range m {
			if !r.Check(v) {
				delete(m, r)
			}
		}
	}

	return m
}

func ruleSliceToMap(rules []Rule) map[Rule]struct{} {
	m := make(map[Rule]struct{})
	for _, rule := range rules {
		m[rule] = struct{}{}
	}
	return m
}

func sumInts(s []int) int {
	sum := 0
	for _, n := range s {
		sum += n
	}
	return sum
}

func checkTickets(tickets []Ticket, rules []Rule) []int {
	var invalid []int
	for _, t := range tickets {
		inv := checkTicket(t, rules)
		invalid = append(invalid, inv...)
	}

	return invalid
}

func filterTickets(tickets []Ticket, rules []Rule) []Ticket {
	var valid []Ticket

	for _, t := range tickets {
		inv := checkTicket(t, rules)
		if len(inv) == 0 {
			valid = append(valid, t)
		}
	}

	return valid
}

func checkTicket(t Ticket, rules []Rule) []int {
	var invalid []int

	for _, v := range t {
		if !checkRule(v, rules) {
			invalid = append(invalid, v)
		}
	}

	return invalid
}

func checkRule(v int, rules []Rule) bool {
	for _, r := range rules {
		if r.Check(v) {
			return true
		}
	}

	return false
}

type Range struct {
	From, To int
}

func (r Range) Check(n int) bool {
	return n >= r.From && n <= r.To
}

type Rule struct {
	Name   string
	Range1 Range
	Range2 Range
}

func (r Rule) Check(n int) bool {
	return r.Range1.Check(n) || r.Range2.Check(n)
}

var matchRule = regexp.MustCompile(`^([\w ]+): (\d+)-(\d+) or (\d+)-(\d+)$`)

func readFile(name string) ([]Rule, Ticket, []Ticket, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	rules, err := readRules(s)
	if err != nil {
		return nil, nil, nil, err
	}

	s.Scan() // discard the "your ticket:" line

	myTicket, err := readTickets(s)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(myTicket) < 1 {
		return nil, nil, nil, errors.New("missing my ticket")
	}

	s.Scan() // discard the "nearby tickets:" line

	tickets, err := readTickets(s)
	if err != nil {
		return nil, nil, nil, err
	}

	return rules, myTicket[0], tickets, nil
}

func readRules(s *bufio.Scanner) ([]Rule, error) {
	var rules []Rule

	for s.Scan() {
		line := s.Text()
		if line == "" {
			return rules, nil
		}

		matches := matchRule.FindStringSubmatch(line)
		if len(matches) != 6 {
			fmt.Printf("Regex konnte nicht matchen: %s\n", line)
			continue
		}

		r := Rule{Name: matches[1]}
		r.Range1.From, _ = strconv.Atoi(matches[2])
		r.Range1.To, _ = strconv.Atoi(matches[3])
		r.Range2.From, _ = strconv.Atoi(matches[4])
		r.Range2.To, _ = strconv.Atoi(matches[5])

		rules = append(rules, r)
	}

	return rules, s.Err()
}

type Ticket []int

func readTickets(s *bufio.Scanner) ([]Ticket, error) {
	var tickets []Ticket

	for s.Scan() {
		line := s.Text()
		if line == "" {
			return tickets, nil
		}

		values := strings.Split(line, ",")
		t := make(Ticket, len(values))
		for i, v := range values {
			t[i], _ = strconv.Atoi(v)
		}

		tickets = append(tickets, t)
	}

	return tickets, s.Err()
}
