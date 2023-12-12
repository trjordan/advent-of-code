package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type galaxy struct {
	galaxyNum int
	row       int
	col       int
}

func galaxyDistance(g1 galaxy, g2 galaxy) int {
	return int(math.Abs(float64(g1.row-g2.row)) + math.Abs(float64(g1.col-g2.col)))
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	galaxies := []galaxy{}
	galaxyNum := 1
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for row, line := range lines {
		for col, point := range line {
			if point == '#' {
				galaxies = append(galaxies, galaxy{
					galaxyNum: galaxyNum,
					row:       row,
					col:       col,
				})
				galaxyNum += 1
			}
		}
	}
	fmt.Println(galaxies)

	// Expand!
	rowHasGalaxy := make([]bool, len(lines))
	colHasGalaxy := make([]bool, len(lines[0]))
	for _, galaxy := range galaxies {
		rowHasGalaxy[galaxy.row] = true
		colHasGalaxy[galaxy.col] = true
	}

	rowOffsets := make([]int, len(rowHasGalaxy))
	colOffsets := make([]int, len(colHasGalaxy))
	extraRows := 0
	for row, offset := range rowHasGalaxy {
		if !offset {
			extraRows += 1
		}
		rowOffsets[row] = extraRows
	}
	extraCols := 0
	for col, offset := range colHasGalaxy {
		if !offset {
			extraCols += 1
		}
		colOffsets[col] = extraCols
	}

	for i := 0; i < len(galaxies); i++ {
		galaxies[i].row += rowOffsets[galaxies[i].row]
		galaxies[i].col += colOffsets[galaxies[i].col]

	}
	fmt.Println("rows", rowOffsets)
	fmt.Println("cols", colOffsets)
	fmt.Println("expanded", galaxies)

	// Distances
	totalDistance := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i; j < len(galaxies); j++ {
			totalDistance += galaxyDistance(galaxies[i], galaxies[j])
		}
	}
	fmt.Println("total", totalDistance)

}
