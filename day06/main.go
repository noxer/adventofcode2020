package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	groups, err := readGroups("input.txt")
	if err != nil {
		fmt.Printf("Konnte die Gruppen nicht einlesen: %s\n", err)
		os.Exit(1)
	}

	sum := 0
	for _, group := range groups {
		answers := findAll(group)
		sum += len(answers)
	}

	fmt.Printf("Anzahl der Antworten: %d\n", sum)
}

func findAny(group []string) string {
	yes := make(map[rune]bool)
	for _, person := range group {
		for _, answer := range person {
			yes[answer] = true
		}
	}

	builder := &strings.Builder{}
	for answer := range yes {
		builder.WriteRune(answer)
	}

	return builder.String()
}

func findAll(group []string) string {
	yes := make(map[rune]int)
	for _, person := range group {
		for _, answer := range person {
			yes[answer]++
		}
	}

	builder := &strings.Builder{}
	for answer, count := range yes {
		if count == len(group) {
			builder.WriteRune(answer)
		}
	}

	return builder.String()
}

func readGroups(name string) ([][]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var groups [][]string
	var group []string

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			if len(group) > 0 {
				groups = append(groups, group)
				group = nil
			}
			continue
		}

		group = append(group, line)
	}

	if len(group) > 0 {
		groups = append(groups, group)
	}

	return groups, nil
}
