package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	movements, err := readMovements("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Bewegungsdaten: %s\n", err)
		os.Exit(1)
	}

	s := &Ship{
		Dir: 90,
	}

	for _, m := range movements {
		m.Move(s)
	}

	fmt.Printf("Lösung 1: %d\n", abs(s.X)+abs(s.Y))

	b := &Position{}
	w := &Position{
		X: 10,
		Y: -1,
	}

	for _, m := range movements {
		m.MoveWaypoint(w, b)
	}

	fmt.Printf("Lösung 2: %d\n", abs(b.X)+abs(b.Y))
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// Ship ...
type Ship struct {
	Dir int
	X   int
	Y   int
}

// Position ...
type Position struct {
	X int
	Y int
}

// Movement ...
type Movement interface {
	Move(*Ship)
	MoveWaypoint(*Position, *Position)
}

// North ...
type North int

// East ...
type East int

// South ...
type South int

// West ...
type West int

// Move ...
func (n North) Move(s *Ship) {
	s.Y -= int(n)
}

// MoveWaypoint ...
func (n North) MoveWaypoint(w, _ *Position) {
	w.Y -= int(n)
}

// Move ...
func (e East) Move(s *Ship) {
	s.X += int(e)
}

// MoveWaypoint ...
func (e East) MoveWaypoint(w, _ *Position) {
	w.X += int(e)
}

// Move ...
func (s South) Move(b *Ship) {
	b.Y += int(s)
}

// MoveWaypoint ...
func (s South) MoveWaypoint(w, _ *Position) {
	w.Y += int(s)
}

// Move ...
func (w West) Move(s *Ship) {
	s.X -= int(w)
}

// MoveWaypoint ...
func (w West) MoveWaypoint(b, _ *Position) {
	b.X -= int(w)
}

// Left ...
type Left int

// Right ...
type Right int

// Forward ...
type Forward int

// Move ...
func (l Left) Move(s *Ship) {
	s.Dir -= int(l)
	if s.Dir < 0 {
		s.Dir += 360
	}
}

// MoveWaypoint ...
func (l Left) MoveWaypoint(w, _ *Position) {
	for ; l > 0; l -= 90 {
		w.Y, w.X = -w.X, w.Y
	}
}

// Move ...
func (r Right) Move(s *Ship) {
	s.Dir += int(r)
	s.Dir = s.Dir % 360
}

// MoveWaypoint ...
func (r Right) MoveWaypoint(w, _ *Position) {
	for ; r > 0; r -= 90 {
		w.Y, w.X = w.X, -w.Y
	}
}

// Move ...
func (f Forward) Move(s *Ship) {
	switch s.Dir {
	case 0:
		s.Y -= int(f)
	case 90:
		s.X += int(f)
	case 180:
		s.Y += int(f)
	case 270:
		s.X -= int(f)
	default:
		panic(fmt.Sprintf("Ungültige Richtung: %d", s.Dir))
	}
}

// MoveWaypoint ...
func (f Forward) MoveWaypoint(w, s *Position) {
	for i := Forward(0); i < f; i++ {
		s.X += w.X
		s.Y += w.Y
	}
}

func readMovements(name string) ([]Movement, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var movements []Movement
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		i, _ := strconv.Atoi(line[1:])
		switch line[0] {
		case 'N':
			movements = append(movements, North(i))
		case 'E':
			movements = append(movements, East(i))
		case 'S':
			movements = append(movements, South(i))
		case 'W':
			movements = append(movements, West(i))
		case 'R':
			movements = append(movements, Right(i))
		case 'L':
			movements = append(movements, Left(i))
		case 'F':
			movements = append(movements, Forward(i))
		}
	}

	return movements, nil
}
