package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	movements, err := readMovements("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Lesen der Bewegungen: %s\n", err)
		os.Exit(1)
	}

	tiles := make(map[Coordinate]struct{})

	for _, movement := range movements {
		state := Coordinate{}

		for _, dir := range movement {
			dir.Move(&state)
		}

		flipTile(tiles, state)
	}

	fmt.Printf("Schwarze Kacheln: %d\n", len(tiles))

	for i := 0; i < 100; i++ {
		tiles = step(tiles)
	}

	fmt.Printf("Schwarze Kacheln nach 100 Iterationen: %d\n", len(tiles))
}

func flipTile(tiles map[Coordinate]struct{}, c Coordinate) {
	if _, ok := tiles[c]; ok {
		delete(tiles, c)
	} else {
		tiles[c] = struct{}{}
	}
}

var directions = []Direction{W, NW, NE, E, SE, SW}

func step(tiles map[Coordinate]struct{}) map[Coordinate]struct{} {
	result := make(map[Coordinate]struct{})
	white := make(map[Coordinate]struct{})

	for tile := range tiles {
		neighbors := 0

		for _, dir := range directions {
			c := dir.Add(tile)
			if _, ok := tiles[c]; ok {
				neighbors++
			} else {
				white[c] = struct{}{}
			}
		}

		if neighbors > 0 && neighbors < 3 {
			result[tile] = struct{}{}
		}
	}

	for tile := range white {
		neighbors := 0

		for _, dir := range directions {
			c := dir.Add(tile)
			if _, ok := tiles[c]; ok {
				neighbors++
			}
		}

		if neighbors == 2 {
			result[tile] = struct{}{}
		}
	}

	return result
}

// Coordinate ...
type Coordinate struct {
	X, Y int
}

// Direction ...
type Direction Coordinate

// Move ...
func (d Direction) Move(c *Coordinate) {
	c.X += d.X
	c.Y += d.Y
}

// Add ...
func (d Direction) Add(c Coordinate) Coordinate {
	c.X += d.X
	c.Y += d.Y

	return c
}

// ...
var (
	W  = Direction{X: -1, Y: 1}
	NW = Direction{X: 0, Y: 1}
	NE = Direction{X: 1, Y: 0}
	E  = Direction{X: 1, Y: -1}
	SE = Direction{X: 0, Y: -1}
	SW = Direction{X: -1, Y: 0}
)

func readMovements(name string) ([][]Direction, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var movements [][]Direction
	for s.Scan() {
		movements = append(movements, readMovement(s.Text()))
	}

	return movements, nil
}

func readMovement(inst string) []Direction {
	prev := rune(0)

	var movement []Direction
	for _, c := range inst {
		switch c {
		case 'e':
			if prev != 0 {
				if prev == 's' {
					movement = append(movement, SE)
				} else {
					movement = append(movement, NE)
				}
				prev = 0
			} else {
				movement = append(movement, E)
			}

		case 'w':
			if prev != 0 {
				if prev == 's' {
					movement = append(movement, SW)
				} else {
					movement = append(movement, NW)
				}
				prev = 0
			} else {
				movement = append(movement, W)
			}

		default:
			prev = c
		}
	}

	return movement
}
