package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type tile struct {
	pipe    rune
	row     int
	col     int
	isStart bool
	north   *tile
	east    *tile
	south   *tile
	west    *tile
}

type direction int

const (
	unknown direction = iota
	north
	east
	south
	west
)

func pr(lbl string, t tile) {
	if t.row == 1 && t.col == 1 {
		fmt.Println(lbl, t)
	}
}

func parseGrid(fullGrid []string) *tile {
	var start *tile
	tileGrid := make([][]tile, len(fullGrid))
	for row, line := range fullGrid {
		tileGrid[row] = make([]tile, len(line))
		for col, pipe := range line {
			tileGrid[row][col] = tile{
				pipe: pipe,
				row:  row,
				col:  col,
			}
			if col > 0 {
				tileGrid[row][col-1].east = &tileGrid[row][col]
				tileGrid[row][col].west = &tileGrid[row][col-1]
			}
			if row > 0 {
				tileGrid[row-1][col].south = &tileGrid[row][col]
				tileGrid[row][col].north = &(tileGrid[row-1][col])
			}

			if tileGrid[row][col].pipe == 'S' {
				start = &tileGrid[row][col]
			}
		}
	}
	return start
}

func neq(t1 tile, t2 tile) bool {
	return !(t1.row == t2.row && t1.col == t2.col)
}

func nextTile(t tile, prev tile) (tile, tile) {
	orig := t.pipe
	if t.north != nil && strings.ContainsRune("|7FS", (*t.north).pipe) && strings.ContainsRune("|JLS", orig) && neq(*t.north, prev) {
		fmt.Println("going north", t.row, t.col)
		return *t.north, t
	}
	if t.east != nil && strings.ContainsRune("-7JS", t.east.pipe) && strings.ContainsRune("-FLS", orig) && neq(*t.east, prev) {
		fmt.Println("going east", t.row, t.col)
		return *t.east, t
	}
	if t.south != nil && strings.ContainsRune("|JLS", (*t.south).pipe) && strings.ContainsRune("|7FS", orig) && neq(*t.south, prev) {
		fmt.Println("going south", t.row, t.col)
		return *t.south, t
	}
	if t.west != nil && strings.ContainsRune("-FLS", (*t.west).pipe) && strings.ContainsRune("-7JS", orig) && neq(*t.west, prev) {
		fmt.Println("going west", t.row, t.col)
		return *t.west, t
	}
	// oh no!!
	panic("no valid tile found")
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	start := *parseGrid(lines)

	// manually start its up!
	prev := start
	next, prev := nextTile(start, start)
	distance := 1
	for next != start {
		next, prev = nextTile(next, prev)
		distance += 1
		fmt.Println(next.pipe, next.row, next.col)
	}
	fmt.Println("Loop finished with distance", distance)
}
