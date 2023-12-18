package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	row, col int
}

func (p Point) next(dir string) Point {
	if dir == "R" {
		return Point{p.row, p.col + 1}
	} else if dir == "L" {
		return Point{p.row, p.col - 1}
	} else if dir == "D" {
		return Point{p.row + 1, p.col}
	} else if dir == "U" {
		return Point{p.row - 1, p.col}
	}
	panic("unknown direction to move")
}

func (p Point) inRange(grid []string) bool {
	return p.row >= 0 && p.row < len(grid) && p.col >= 0 && p.row < len(grid[0])
}

func printGrid(grid map[Point]bool) {
	minRow := 0
	minCol := 0
	maxRow := 0
	maxCol := 0
	for p, _ := range grid {
		if p.row > maxRow {
			maxRow = p.row
		}
		if p.row < minRow {
			minRow = p.row
		}
		if p.col > maxCol {
			maxCol = p.col
		}
		if p.col < minCol {
			minCol = minCol
		}
	}

	// minRow = -50

	// maxRow = 25
	// maxCol = 100

	for i := minRow; i < maxRow+1; i++ {
		fmt.Printf("%03d ", i)
		for j := minCol; j < maxCol+1; j++ {
			if grid[Point{i, j}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//
	grid := map[Point]bool{}
	curPoint := Point{0, 0}
	grid[curPoint] = true
	for _, line := range lines {
		fields := strings.Fields(line)
		dir := fields[0]
		numSteps, _ := strconv.Atoi(fields[1])
		for i := 0; i < numSteps; i++ {
			curPoint = curPoint.next(dir)
			grid[curPoint] = true
		}

	}

	fmt.Println("unfilled")
	printGrid(grid)

	// Flood fill
	flooded := map[Point]bool{} // Make a copy so we can look for walls as we fill
	for p, _ := range grid {
		flooded[p] = true
	}
	openSet := []Point{Point{1, 1}} // Seems generically safe?
	flooded[openSet[0]] = true
	for len(openSet) > 0 {
		// Pop the current point
		curPoint := openSet[0]
		openSet = openSet[1:]
		fmt.Println("working on ", curPoint, len(openSet))

		// Look around this point and enqueue and non-wall points
		candidates := []Point{
			curPoint.next("L"),
			curPoint.next("R"),
			curPoint.next("U"),
			curPoint.next("D"),
		}
		fmt.Println(candidates)
		for _, c := range candidates {
			if !grid[c] && !flooded[c] {
				// Inside point, not seen before, so fill it and queue up exploration
				fmt.Println("queueing", c, !grid[c], !flooded[c])
				openSet = append(openSet, c)
				flooded[c] = true
			}
		}
	}

	fmt.Println("unfilled")
	printGrid(grid)
	fmt.Println("filled")
	printGrid(flooded)
	fmt.Println("Num filled", len(flooded))

}
