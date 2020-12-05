package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	bps, err := readBoardingPasses("input.txt")
	if err != nil {
		fmt.Printf("Konnte die Boarding Pässe nicht einlesen: %s\n", err)
	}

	plane := [8][128]bool{}

	for _, bp := range bps {
		plane[bp.Column][bp.Row] = true
	}

out:
	for c, col := range plane {
		for r, occupied := range col[1 : len(col)-1] {
			if !occupied && plane[c][r+2] && plane[c][r] {
				fmt.Printf("Unser Sitz ist: Reihe %d, Sitz %d\n", r+1, c)
				id := BoardingPass{Row: r + 1, Column: c}.SeatID()
				fmt.Printf("Unsere Sitz ID ist %d\n", id)
				break out
			}
		}
	}
}

// BoardingPass ...
type BoardingPass struct {
	SeatNumber string
	Row        int
	Column     int
}

// SeatID ...
func (b BoardingPass) SeatID() int {
	return b.Row*8 + b.Column
}

func readBoardingPasses(name string) ([]BoardingPass, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var boardingPasses []BoardingPass
	for s.Scan() {
		b := BoardingPass{SeatNumber: s.Text()}
		b.Row, b.Column = calculateRowColumn(s.Text())
		boardingPasses = append(boardingPasses, b)
	}

	return boardingPasses, nil
}

func calculateRowColumn(seatNum string) (row, col int) {
	row = interpretBinary(seatNum[:7], 'B')
	col = interpretBinary(seatNum[7:], 'R')
	return
}

func interpretBinary(num string, one rune) int {
	i := 0
	for _, c := range num {
		i <<= 1
		if c == one {
			i |= 1
		}
	}
	return i
}
