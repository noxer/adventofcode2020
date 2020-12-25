package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	tiles, err := readTiles("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Kacheln: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d Kacheln geladen...\n", len(tiles))

	db := NewTileDB(tiles)

	var corner Tile
	for _, tile := range tiles {
		neighbors := 0
		for _, edge := range tile.Edges() {
			if db.findByEdge(edge, tile.ID) != nil {
				neighbors++
			}
		}

		if neighbors == 2 {
			corner = tile
			break
		}
	}

	for {
		RotateRight(&corner)

		if db.findByEdge(corner.West(), corner.ID) != nil {
			continue
		}
		if db.findByEdge(corner.North(), corner.ID) != nil {
			continue
		}

		break
	}

	m := map[Coord]*Tile{
		Coord{0, 0}: &corner,
	}

out:
	for y := 0; ; y++ {
		for x := 0; ; x++ {
			tile := m[Coord{x, y}]
			if tile == nil {
				break out
			}

			rightTile := db.findByEdge(tile.East(), tile.ID)
			if rightTile != nil {
				turnTileToMatch(rightTile, tile.East(), West)
				m[Coord{x + 1, y}] = rightTile
			}

			bottomTile := db.findByEdge(tile.South(), tile.ID)
			if bottomTile != nil {
				turnTileToMatch(bottomTile, tile.South(), North)
				m[Coord{x, y + 1}] = bottomTile
			}

			if rightTile == nil {
				break
			}
		}
	}

	completeMap := [96][96]byte{}
	for y, row := range completeMap {
		for x := range row {
			completeMap[y][x] = getFromMap(m, x, y)
		}
	}

	for i := 0; i < 4; i++ {
		for y := 0; y < 96; y++ {
			for x := 0; x < 96; x++ {
				if matchSeamonster(completeMap, x, y) {
					fmt.Printf("Monster bei (%d, %d)\n", x, y)
				}
			}
		}

		completeMap = RotateMap(completeMap)
	}

	completeMap = MirrorMap(completeMap)

	for i := 0; i < 4; i++ {
		var monsters []Coord

		for y := 0; y < 96; y++ {
			for x := 0; x < 96; x++ {
				if matchSeamonster(completeMap, x, y) {
					fmt.Printf("Monster bei (%d, %d)\n", x, y)
					monsters = append(monsters, Coord{x, y})
				}
			}
		}

		if len(monsters) > 0 {
			fmt.Printf("Wir haben %d Monster gefunden...\n", len(monsters))
			for _, monster := range monsters {
				completeMap = deleteMonster(completeMap, monster.X, monster.Y)
			}

			waveCount := 0
			for _, row := range completeMap {
				for _, b := range row {
					if b == '#' {
						waveCount++
					}
				}
			}

			fmt.Printf("%d Wellen sind in dem Gebiet\n", waveCount)
		}

		completeMap = RotateMap(completeMap)
	}

}

func matchSeamonster(m [96][96]byte, x, y int) bool {
	monster := [][]byte{
		[]byte("                  # "),
		[]byte("#    ##    ##    ###"),
		[]byte(" #  #  #  #  #  #   "),
	}

	for monsterY, row := range monster {
		for monsterX, b := range row {
			if b != '#' {
				continue
			}

			if y+monsterY >= 96 || x+monsterX >= 96 {
				return false
			}

			if m[y+monsterY][x+monsterX] != '#' {
				return false
			}
		}
	}

	return true
}

func deleteMonster(m [96][96]byte, x, y int) [96][96]byte {
	monster := [][]byte{
		[]byte("                  # "),
		[]byte("#    ##    ##    ###"),
		[]byte(" #  #  #  #  #  #   "),
	}

	for monsterY, row := range monster {
		for monsterX, b := range row {
			if b != '#' {
				continue
			}

			m[y+monsterY][x+monsterX] = ' '
		}
	}

	return m
}

type Coord struct {
	X, Y int
}

// Tile ...
type Tile struct {
	ID   int
	Data [10][10]byte
}

func (t Tile) North() Edge {
	return Edge(t.Data[0])
}

func (t Tile) South() Edge {
	return Edge(t.Data[9])
}

func (t Tile) West() Edge {
	var e Edge
	for i, row := range t.Data {
		e[i] = row[0]
	}

	return e
}

func (t Tile) East() Edge {
	var e Edge
	for i, row := range t.Data {
		e[i] = row[9]
	}

	return e
}

func (t Tile) Edges() [4]Edge {
	var edges [4]Edge
	edges[0] = t.North()
	edges[1] = t.East()
	edges[2] = t.South()
	edges[3] = t.West()
	return edges
}

type TileDB struct {
	tiles map[int]*Tile
	index map[Edge][]int
}

func NewTileDB(tiles []Tile) *TileDB {
	db := &TileDB{
		tiles: make(map[int]*Tile),
		index: make(map[Edge][]int),
	}

	for i, tile := range tiles {
		db.tiles[tile.ID] = &tiles[i]

		db.index[tile.East()] = append(db.index[tile.East()], tile.ID)
		db.index[tile.South()] = append(db.index[tile.South()], tile.ID)
		db.index[tile.West()] = append(db.index[tile.West()], tile.ID)
		db.index[tile.North()] = append(db.index[tile.North()], tile.ID)

		db.index[tile.East().Invert()] = append(db.index[tile.East().Invert()], tile.ID)
		db.index[tile.South().Invert()] = append(db.index[tile.South().Invert()], tile.ID)
		db.index[tile.West().Invert()] = append(db.index[tile.West().Invert()], tile.ID)
		db.index[tile.North().Invert()] = append(db.index[tile.North().Invert()], tile.ID)
	}

	return db
}

func (db *TileDB) findByEdge(edge Edge, ignoreID int) *Tile {
	ids, ok := db.index[edge]
	if ok {
		for _, id := range ids {
			if id == ignoreID {
				continue
			}

			return db.tiles[id]
		}
	}

	return nil
}

func readTiles(name string) ([]Tile, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var tiles []Tile
	for s.Scan() {
		t, err := readTile(s)
		if err != nil {
			return tiles, err
		}

		tiles = append(tiles, t)
	}

	return tiles, nil
}

func readTile(s *bufio.Scanner) (Tile, error) {
	s.Scan()
	label := strings.TrimSuffix(strings.TrimPrefix(s.Text(), "Tile "), ":")

	id, err := strconv.Atoi(label)
	if err != nil {
		return Tile{}, err
	}

	t := Tile{
		ID: id,
	}

	for i := range t.Data {
		s.Scan()
		copy(t.Data[i][:], s.Bytes())
	}

	return t, nil
}

// Edge ...
type Edge [10]byte

func (e Edge) Print() {
	fmt.Println(string(e[:]))
}

// Equals ...
func (e Edge) Equals(o Edge) int {
	for i, b := range e {
		if b != o[i] {
			return e.equalsInv(o)
		}
	}

	return 1
}

func (e Edge) equalsInv(o Edge) int {
	for i, b := range e {
		if b != o[9-i] {
			return 0
		}
	}

	return 2
}

func (e Edge) Invert() Edge {
	var r Edge
	for i, b := range e {
		r[9-i] = b
	}
	return r
}

const (
	North = iota
	East
	South
	West
)

func turnTileToMatch(tile *Tile, edge Edge, side int) {
	matchingSide := -1
	for i, e := range tile.Edges() {
		if e.Equals(edge) != 0 {
			matchingSide = i
			break
		}
	}

	for i := 0; i < (side-matchingSide+4)%4; i++ {
		RotateRight(tile)
	}

	matchingEdge := tile.Edges()[side]

	switch matchingEdge.Equals(edge) {
	case 0:
		panic(fmt.Sprintf("Fehler beim Drehen, Kanten passen nicht nach Drehung: TileID %d", tile.ID))

	case 1:
		// the edges line up exactly, we are done
		return

	case 2:
		if side == 0 || side == 2 {
			MirrorV(tile)
		} else {
			MirrorH(tile)
		}
	}
}

// RotateRight ...
func RotateRight(tile *Tile) {
	N := 10
	a := tile.Data

	// Traverse each cycle
	for i := 0; i < N/2; i++ {
		for j := i; j < N-i-1; j++ {

			// Swap elements of each cycle
			// in clockwise direction
			temp := a[i][j]
			a[i][j] = a[N-1-j][i]
			a[N-1-j][i] = a[N-1-i][N-1-j]
			a[N-1-i][N-1-j] = a[j][N-1-i]
			a[j][N-1-i] = temp
		}
	}

	tile.Data = a
}

// RotateMap ...
func RotateMap(m [96][96]byte) [96][96]byte {
	N := 96
	a := m

	// Traverse each cycle
	for i := 0; i < N/2; i++ {
		for j := i; j < N-i-1; j++ {

			// Swap elements of each cycle
			// in clockwise direction
			temp := a[i][j]
			a[i][j] = a[N-1-j][i]
			a[N-1-j][i] = a[N-1-i][N-1-j]
			a[N-1-i][N-1-j] = a[j][N-1-i]
			a[j][N-1-i] = temp
		}
	}

	return a
}

// MirrorV ...
func MirrorV(tile *Tile) {
	for _, row := range tile.Data {
		invertRow(row[:])
	}
}

// MirrorH ...
func MirrorH(tile *Tile) {
	for i := 0; i < len(tile.Data)/2; i++ {
		tile.Data[i], tile.Data[len(tile.Data)-i-1] = tile.Data[len(tile.Data)-i-1], tile.Data[i]
	}
}

// MirrorH ...
func MirrorMap(m [96][96]byte) [96][96]byte {
	for i := 0; i < len(m)/2; i++ {
		m[i], m[len(m)-i-1] = m[len(m)-i-1], m[i]
	}
	return m
}

func invertRow(row []byte) {
	for i := 0; i < len(row)/2; i++ {
		row[i], row[len(row)-i-1] = row[len(row)-i-1], row[i]
	}
}

// Print ...
func (t Tile) Print() {
	fmt.Printf("Tile %d:\n", t.ID)
	for _, row := range t.Data {
		fmt.Println(string(row[:]))
	}
}

func getFromMap(m map[Coord]*Tile, x, y int) byte {
	tileX := x / 8
	tileY := y / 8

	tile := m[Coord{tileX, tileY}]
	if tile == nil {
		return 0
	}

	x %= 8
	y %= 8

	return tile.Data[y+1][x+1]
}
