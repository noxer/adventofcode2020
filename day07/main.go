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
	rules, err := readRules("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Regeln: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Anzahl der Regeln: %d\n", len(rules))

	inv := invertTree(rules)
	count := countParents(inv, "shiny gold")
	fmt.Printf("Farben, die unsere Farbe beinhalten können: %d\n", count)

	n := countContained(rules, "shiny gold") - 1
	fmt.Printf("Unsere Taschen muss %d andere Taschen enthalten...\n", n)
}

func invertTree(rules map[string]Rule) map[string]map[string]struct{} {
	m := make(map[string]map[string]struct{})

	for color, rule := range rules {
		for contained := range rule.CanContain {
			existing := m[contained]
			if existing == nil {
				existing = make(map[string]struct{})
			}

			existing[color] = struct{}{}
			m[contained] = existing
		}
	}

	return m
}

func countContained(in map[string]Rule, color string) int {
	n := 1
	for contained, count := range in[color].CanContain {
		n += countContained(in, contained) * count
	}
	return n
}

func countParents(in map[string]map[string]struct{}, color string) int {
	seen := map[string]struct{}{
		color: struct{}{},
	}

	findParents(in, color, seen)
	return len(seen) - 1
}

func findParents(in map[string]map[string]struct{}, color string, seen map[string]struct{}) {
	for parent := range in[color] {
		if _, ok := seen[parent]; ok {
			continue
		}

		seen[parent] = struct{}{}
		findParents(in, parent, seen)
	}
}

// Rule ...
type Rule struct {
	Color      string
	CanContain map[string]int
}

var (
	matchLine = regexp.MustCompile(`^(\w+ \w+) bags contain (.+).$`)
	matchPart = regexp.MustCompile(`^(\d+) (\w+ \w+) bags?$`)
)

// ParseRule ...
func ParseRule(line string) (Rule, error) {
	groups := matchLine.FindStringSubmatch(line)
	if len(groups) != 3 {
		return Rule{}, errors.New("could not match")
	}

	r := Rule{
		Color:      groups[1],
		CanContain: make(map[string]int),
	}

	if groups[2] == "no other bags" {
		return r, nil
	}

	parts := strings.Split(groups[2], ", ")
	for _, part := range parts {
		segments := matchPart.FindStringSubmatch(part)
		if len(segments) != 3 {
			fmt.Printf("Konnte Segment %s nicht parsen!\n", part)
			continue
		}

		n, _ := strconv.Atoi(segments[1])
		r.CanContain[segments[2]] += n
	}

	return r, nil
}

func readRules(name string) (map[string]Rule, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	rules := make(map[string]Rule)
	for s.Scan() {
		r, err := ParseRule(s.Text())
		if err != nil {
			continue
		}

		rules[r.Color] = r
	}

	return rules, nil
}
