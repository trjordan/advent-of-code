package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printGrid(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func rotate(lines []string) []string {
	ret := make([]string, len(lines[0]))
	for i := len(lines) - 1; i >= 0; i-- {
		for j := 0; j < len(lines[0]); j++ {
			ret[j] += string(lines[i][j])
		}
	}
	return ret
}

func transpose(lines []string) []string {
	ret := make([]string, len(lines[0]))
	for i := 0; i < len(lines); i++ {
		for j, chr := range lines[i] {
			ret[j] += string(chr)
		}
	}
	return ret
}

// Rolls a single row west (right to left)
func rollRow(line string) string {
	// Make a list of stops, numBoulders
	stops := [][2]int{}
	curStop := -1
	curBoulders := 0
	for i, chr := range line {
		if chr == '#' {
			stops = append(stops, [2]int{curStop, curBoulders})
			curStop = i
			curBoulders = 0
		} else if chr == 'O' {
			curBoulders += 1
		}
	}
	stops = append(stops, [2]int{curStop, curBoulders})

	// Recreate the string
	newLine := ""
	remainingBoulders := stops[0][1]
	stops = stops[1:]
	for i := 0; i < len(line); i++ {
		if len(stops) > 0 && stops[0][0] == i {
			// On to the next pair
			remainingBoulders = stops[0][1]
			stops = stops[1:]
			newLine += "#"
		} else if remainingBoulders > 0 {
			// Place a boulder
			newLine += "O"
			remainingBoulders -= 1
		} else {
			// place a non-boulder
			newLine += "."
		}
	}

	return newLine
}

// Roll all the boulders north
func rollWest(grid []string) []string {
	newGrid := []string{}
	for _, line := range grid {
		newLine := rollRow(line)
		newGrid = append(newGrid, newLine)
	}
	return newGrid
}

func multiRotate(grid []string, numRotations int, doRoll bool) []string {
	for i := 0; i < numRotations; i++ {
		grid = rotate(grid)
		if doRoll {
			grid = rollWest(grid)
		}
	}
	return grid
}

// Walk through and get the total load
func sumWeight(grid []string) int {
	totalWeight := 0
	for i, line := range grid {
		weight := len(grid) - i
		for _, chr := range line {
			if chr == 'O' {
				totalWeight += weight
			}
		}
	}
	return totalWeight
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	grid := []string{}
	for scanner.Scan() {
		grid = append(grid, scanner.Text())

	}
	// Prep since west is the easy way to roll
	printGrid(grid)
	grid = multiRotate(grid, 2, false)
	cycleMap := map[string]int{}       // flattened grid -> step num
	targetIterations := 1000000000 - 1 // stupid puzzle counts from 1
	for i := 0; i < targetIterations; i++ {
		grid = multiRotate(grid, 4, true)

		// Manual check
		// fmt.Println("after ", i, sumWeight(multiRotate(grid, 2, false)))
		// printGrid(multiRotate(grid, 2, false))

		// Check for cycles
		cycleKey := strings.Join(multiRotate(grid, 2, false), "")
		prevIteration, ok := cycleMap[cycleKey]
		if ok {
			// We found a cycle! Fast-forward to the next partial cycle, to
			// match what it would be at targetIterations
			cycleLength := i - prevIteration
			offset := (targetIterations - i) % cycleLength
			fmt.Println("Found cycle ", cycleLength, i, offset)
			grid = multiRotate(grid, offset*4, true)
			break
		}
		// Store the key, keep going
		cycleMap[cycleKey] = i
	}

	grid = multiRotate(grid, 2, false)
	total := sumWeight(grid)
	fmt.Println("final grid")
	printGrid(grid)
	fmt.Println("total", total)

}
