package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type tile struct {
	pipe     rune
	row      int
	col      int
	isStart  bool
	isOnPath bool
	north    *tile
	east     *tile
	south    *tile
	west     *tile
}

func pr(lbl string, t tile) {
	if t.row == 1 && t.col == 1 {
		fmt.Println(lbl, t)
	}
}

func parseGrid(fullGrid []string) ([][]tile, *tile) {
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
	return tileGrid, start
}

func neq(t1 tile, t2 tile) bool {
	return !(t1.row == t2.row && t1.col == t2.col)
}

func nextTile(t *tile, prev *tile) (*tile, *tile) {
	orig := t.pipe
	t.isOnPath = true
	if t.north != nil && strings.ContainsRune("|7FS", t.north.pipe) && strings.ContainsRune("|JLS", orig) && neq(*t.north, *prev) {
		fmt.Println("going north", t.row, t.col)
		return t.north, t
	}
	if t.east != nil && strings.ContainsRune("-7JS", t.east.pipe) && strings.ContainsRune("-FLS", orig) && neq(*t.east, *prev) {
		fmt.Println("going east", t.row, t.col)
		return t.east, t
	}
	if t.south != nil && strings.ContainsRune("|JLS", t.south.pipe) && strings.ContainsRune("|7FS", orig) && neq(*t.south, *prev) {
		fmt.Println("going south", t.row, t.col)
		return t.south, t
	}
	if t.west != nil && strings.ContainsRune("-FLS", t.west.pipe) && strings.ContainsRune("-7JS", orig) && neq(*t.west, *prev) {
		fmt.Println("going west", t.row, t.col)
		return t.west, t
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
	tileGrid, start := parseGrid(lines)

	// manually start its up! Collect the coords of pipes
	prev := start
	next, prev := nextTile(start, start)
	distance := 1
	for next != start {
		next, prev = nextTile(next, prev)
		distance += 1
		fmt.Println(next.pipe, next.row, next.col)
	}
	fmt.Println("grid", tileGrid)
	fmt.Println("distance", distance)

	// Cross-cut lines
	// - empty spots after odd numbers of crossings are inside
	// - empty sports after even numbers of crossings are inside
	totalInside := 0
	for row := 0; row < len(tileGrid); row++ {
		tileLine := tileGrid[row]
		fmt.Println(tileLine)
		isInside := false
		lastTraverseEntry := ' '
		for col := 0; col < len(tileGrid[row]); col++ {
			t := tileGrid[row][col]
			if t.isOnPath {
				if t.pipe == 'F' || t.pipe == 'L' {
					// We're about to traverse a pipe, capture how we got in
					lastTraverseEntry = t.pipe
				} else if (lastTraverseEntry == 'L' && t.pipe == '7') || (lastTraverseEntry == 'F' && t.pipe == 'J') {
					// We passed through a traverse that flipped us
					isInside = isInside != true
				} else if t.pipe == 'J' || t.pipe == '7' {
					// We passed through a traverse that didn't flip us
					lastTraverseEntry = ' '
				} else if t.pipe == '|' {
					// Simple flip
					isInside = isInside != true
				}
			} else if isInside {
				// Count it
				totalInside += 1
			}
		}
	}
	fmt.Println("total inside", totalInside)
}
