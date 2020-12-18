package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	m, err := readMap("input.txt")
	if err != nil {
		fmt.Printf("Konnte die Map nicht einlesen: %s\n", err)
		os.Exit(1)
	}

	for i := 0; i < 6; i++ {
		m = step4D(m)
	}

	fmt.Printf("%d Würfel sind aktiv!\n", len(m))
}

func step(m map[Coordinate]struct{}) map[Coordinate]struct{} {
	next := make(map[Coordinate]struct{})
	unset := make(map[Coordinate]struct{})

	// check if cube becomes inactive
	for c := range m {
		n := countNeighbors(m, c, unset)
		if n == 2 || n == 3 {
			next[c] = struct{}{}
		}
	}

	// check if cube becomes active
	for u := range unset {
		n := countNeighbors(m, u, nil)
		if n == 3 {
			next[u] = struct{}{}
		}
	}

	return next
}

func step4D(m map[Coordinate]struct{}) map[Coordinate]struct{} {
	next := make(map[Coordinate]struct{})
	unset := make(map[Coordinate]struct{})

	// check if cube becomes inactive
	for c := range m {
		n := countNeighbors4D(m, c, unset)
		if n == 2 || n == 3 {
			next[c] = struct{}{}
		}
	}

	// check if cube becomes active
	for u := range unset {
		n := countNeighbors4D(m, u, nil)
		if n == 3 {
			next[u] = struct{}{}
		}
	}

	return next
}

func countNeighbors(m map[Coordinate]struct{}, c Coordinate, unset map[Coordinate]struct{}) int {
	neighbors := 0

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				if x == 0 && y == 0 && z == 0 {
					continue
				}

				n := Coordinate{
					X: c.X + x,
					Y: c.Y + y,
					Z: c.Z + z,
				}
				if _, ok := m[n]; ok {
					neighbors++
				} else {
					if unset != nil {
						unset[n] = struct{}{}
					}
				}
			}
		}
	}

	return neighbors
}

func countNeighbors4D(m map[Coordinate]struct{}, c Coordinate, unset map[Coordinate]struct{}) int {
	neighbors := 0

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				for w := -1; w <= 1; w++ {
					if x == 0 && y == 0 && z == 0 && w == 0 {
						continue
					}

					n := Coordinate{
						X: c.X + x,
						Y: c.Y + y,
						Z: c.Z + z,
						W: c.W + w,
					}
					if _, ok := m[n]; ok {
						neighbors++
					} else {
						if unset != nil {
							unset[n] = struct{}{}
						}
					}
				}
			}
		}
	}

	return neighbors
}

// Coordinate ...
type Coordinate struct {
	X, Y, Z, W int
}

func readMap(name string) (map[Coordinate]struct{}, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	m := make(map[Coordinate]struct{})
	c := Coordinate{}
	for s.Scan() {
		for _, b := range s.Bytes() {
			if b == '#' {
				m[c] = struct{}{}
			}
			c.X++
		}
		c.X = 0
		c.Y++
	}

	return m, nil
}
