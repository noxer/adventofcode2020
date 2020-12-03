package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	m, err := loadMap("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Map: %s\n", err)
		os.Exit(1)
	}

	slopes := []Slope{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}}

	prod := 1
	for _, s := range slopes {
		trees := checkSlope(m, s.Right, s.Down)
		fmt.Printf("Teste Slope %d, %d -> %d Bäume getroffen\n", s.Right, s.Down, trees)

		prod *= trees
	}

	fmt.Printf("Lösung: %d\n", prod)
}

// Slope ...
type Slope struct {
	Right int
	Down  int
}

func checkSlope(m *Map, right, down int) int {
	x, y := 0, 0
	trees := 0
	for !m.IsBeyondMap(x, y) {
		// move the sled
		x += right
		y += down

		if m.Tree(x, y) {
			trees++
		}
	}

	return trees
}

// Map ...
type Map struct {
	Data   []byte
	Width  int
	Height int
}

// Get ...
func (m *Map) Get(x, y int) byte {
	if y >= m.Height || y < 0 {
		return 0
	}

	x = x % m.Width

	return m.Data[y*m.Width+x]
}

// Tree ...
func (m *Map) Tree(x, y int) bool {
	b := m.Get(x, y)
	return b == '#'
}

// IsBeyondMap ...
func (m *Map) IsBeyondMap(x, y int) bool {
	return y >= len(m.Data) || y < 0
}

func (m *Map) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "Map %dx%d\n", m.Width, m.Height)

	for l := 0; l < m.Height; l++ {
		b.Write(m.Data[l*m.Width : l*m.Width+m.Width])
		b.WriteByte('\n')
	}

	return b.String()
}

func loadMap(name string) (*Map, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var data []byte
	width := 0
	height := 0
	for s.Scan() {
		p := s.Bytes()
		if len(p) == 0 {
			continue
		}

		if len(p) > width {
			width = len(p)
		}

		data = append(data, p...)
		height++
	}

	return &Map{
		Data:   data,
		Width:  width,
		Height: height,
	}, nil
}
