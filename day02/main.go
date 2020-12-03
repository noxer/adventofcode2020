package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var matchLine = regexp.MustCompile(`^(\d+)-(\d+) (\w): (\w+)$`)

func main() {
	passwords, err := readPasswords("input.txt")
	if err != nil {
		fmt.Printf("Konnte die Eingabeliste nicht lesen: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d Passwörter gefunden!\n", len(passwords))

	n := 0
	for _, password := range passwords {
		if password.Validate2() {
			n++
		}
	}

	fmt.Printf("%d gültige Passwörter gefunden!\n", n)
}

// Password ...
type Password struct {
	Min  int
	Max  int
	Char string
	Pass string
}

// Validate ...
func (p Password) Validate() bool {
	n := strings.Count(p.Pass, p.Char)
	return n >= p.Min && n <= p.Max
}

// Validate2 ...
func (p Password) Validate2() bool {
	a, b := p.Pass[p.Min-1], p.Pass[p.Max-1]
	c := p.Char[0]

	if a != c && b != c {
		return false
	}

	return a != b
}

func readPasswords(name string) ([]Password, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var passwords []Password
	for s.Scan() {
		match := matchLine.FindStringSubmatch(s.Text())

		password := Password{}
		password.Min, _ = strconv.Atoi(match[1])
		password.Max, _ = strconv.Atoi(match[2])
		password.Char = match[3]
		password.Pass = match[4]

		passwords = append(passwords, password)
	}

	return passwords, nil
}
