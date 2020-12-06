package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	passports, err := readPassports("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Reisepässe: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d Pässe gefunden.\n", len(passports))

	valid := 0
	for _, passport := range passports {
		if passport.ValidStrict() {
			valid++
		}
	}

	fmt.Printf("Erstes: %#v\n", passports[0])
	fmt.Printf("Letztes: %#v\n", passports[len(passports)-1])

	fmt.Printf("%d valide Pässe gefunden.\n", valid)
}

// Passport ...
type Passport map[string]string

var requiredFields = []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}

// Valid ...
func (p Passport) Valid() bool {
	for _, field := range requiredFields {
		if _, ok := p[field]; !ok {
			return false
		}
	}

	return true
}

var (
	matchHeight = regexp.MustCompile(`^(\d+)(in|cm)$`)
	matchHair   = regexp.MustCompile(`^#[0-9a-f]{6}$`)
	matchEye    = regexp.MustCompile(`^amb|blu|brn|gry|grn|hzl|oth$`)
	matchPID    = regexp.MustCompile(`^\d{9}$`)
)

// ValidStrict ...
func (p Passport) ValidStrict() bool {
	return validateNum(p["byr"], 1920, 2002) &&
		validateNum(p["iyr"], 2010, 2020) &&
		validateNum(p["eyr"], 2020, 2030) &&
		validateHeight(p["hgt"]) &&
		matchHair.MatchString(p["hcl"]) &&
		matchEye.MatchString(p["ecl"]) &&
		matchPID.MatchString(p["pid"])
}

func validateNum(num string, lower, upper int) bool {
	n, err := strconv.Atoi(num)
	if err != nil {
		return false
	}

	return n >= lower && n <= upper
}

func validateHeight(value string) bool {
	parts := matchHeight.FindStringSubmatch(value)
	if len(parts) != 3 {
		return false
	}

	n, _ := strconv.Atoi(parts[1])
	if parts[2] == "cm" {
		return n >= 150 && n <= 193
	}

	return n >= 59 && n <= 76
}

func readPassports(name string) ([]Passport, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(scanPassport)

	var passports []Passport
	for s.Scan() {
		raw := bytes.ReplaceAll(s.Bytes(), []byte{' '}, []byte{'\n'})
		parts := bytes.Split(raw, []byte{'\n'})

		p := make(Passport, len(parts))
		for _, part := range parts {
			kv := bytes.SplitN(part, []byte{':'}, 2)
			if len(kv) != 2 {
				continue
			}

			key := string(bytes.TrimSpace(kv[0]))
			p[key] = string(bytes.TrimSpace(kv[1]))
		}

		passports = append(passports, p)
	}

	return passports, nil
}

var passportSeparator = []byte("\n\n")

func scanPassport(data []byte, atEOF bool) (advance int, token []byte, err error) {
	wasNewline := false
	index := -1
	for i, b := range data {
		if b == '\n' {
			if wasNewline {
				index = i
				break
			}
			wasNewline = true
			continue
		}

		if b != '\r' {
			wasNewline = false
		}
	}

	if atEOF && index < 0 {
		index = len(data) - 1
	}
	if index < 0 {
		return 0, nil, nil
	}

	return index + 1, bytes.TrimSpace(data[:index]), nil
}
