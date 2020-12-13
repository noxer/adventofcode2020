package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	area, err := readSeats("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen des Sitzplans: %s\n", err)
		os.Exit(1)
	}

	var changed int
	for area, changed = area.advance(); changed != 0; area, changed = area.advance() {
		//		fmt.Printf("Geändert: %d\nBenutzt: %d\n\n", changed, area.countOccupied())
	}

	fmt.Printf("Lösung: %d\n", area.countOccupied())
}

type Area struct {
	Width  int
	Height int
	Data   [][]byte
}

func (a *Area) Print() {
	for _, r := range a.Data {
		fmt.Println(string(r))
	}
}

func (a *Area) advance() (*Area, int) {
	c := createCopy(a)

	changes := 0
	for y := 0; y < a.Height; y++ {
		for x := 0; x < a.Width; x++ {
			occupied := a.countAdjacent(x, y)

			switch a.Data[y][x] {
			case '#':
				if occupied >= 4 {
					c.Data[y][x] = 'L'
					changes++
				} else {
					c.Data[y][x] = '#'
				}
			case 'L':
				if occupied == 0 {
					c.Data[y][x] = '#'
					changes++
				} else {
					c.Data[y][x] = 'L'
				}
			}
		}
	}

	return c, changes
}

func (a *Area) advanceLOS() (*Area, int) {
	c := createCopy(a)

	changes := 0
	for y := 0; y < a.Height; y++ {
		for x := 0; x < a.Width; x++ {
			occupied := a.countInLineOfSight(x, y)

			switch a.Data[y][x] {
			case '#':
				if occupied >= 5 {
					c.Data[y][x] = 'L'
					changes++
				} else {
					c.Data[y][x] = '#'
				}
			case 'L':
				if occupied == 0 {
					c.Data[y][x] = '#'
					changes++
				} else {
					c.Data[y][x] = 'L'
				}
			}
		}
	}

	return c, changes
}

func (a *Area) countOccupied() int {
	sum := 0
	for _, r := range a.Data {
		for _, b := range r {
			if b == '#' {
				sum++
			}
		}
	}
	return sum
}

func (a *Area) countInLineOfSight(x, y int) int {
	sum := 0
	for v := -1; v <= 1; v++ {
		for k := -1; k <= 1; k++ {
			if k == 0 && v == 0 {
				continue
			}

			if a.getOccupiedInLine(x, y, v, k) {
				sum++
			}
		}
	}

	return sum
}

func (a *Area) countAdjacent(x, y int) int {
	sum := 0
	for v := -1; v <= 1; v++ {
		for k := -1; k <= 1; k++ {
			if k == 0 && v == 0 {
				continue
			}

			if a.getOccupied(x+v, y+k) {
				sum++
			}
		}
	}

	return sum
}

func createCopy(a *Area) *Area {
	area := make([][]byte, a.Height)
	for i := range area {
		area[i] = bytes.Repeat([]byte{'.'}, a.Width)
	}
	return &Area{
		Width:  a.Width,
		Height: a.Height,
		Data:   area,
	}
}

func readSeats(name string) (*Area, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var area [][]byte
	for s.Scan() {
		c := make([]byte, len(s.Bytes()))
		copy(c, s.Bytes())
		area = append(area, bytes.TrimSpace(c))
	}

	return &Area{
		Width:  len(area[0]),
		Height: len(area),
		Data:   area,
	}, s.Err()
}

func (a *Area) getOccupied(x, y int) bool {
	if x < 0 || x >= a.Width || y < 0 || y >= a.Height {
		return false
	}

	if a.Data[y][x] == '#' {
		return true
	}

	return false
}

func (a *Area) getOccupiedInLine(sx, sy, dx, dy int) bool {
	for x, y := sx+dx, sy+dy; x >= 0 && x < a.Width && y >= 0 && y < a.Height; x, y = x+dx, y+dy {
		switch a.Data[y][x] {
		case 'L':
			return false
		case '#':
			return true
		}
	}

	return false
}
