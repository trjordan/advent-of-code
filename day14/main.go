package main

import (
	"bufio"
	"fmt"
	"os"
)

func printGrid(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
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
	fmt.Println("orig  ", line)
	fmt.Println("rolled", newLine)

	return newLine
}

// Roll all the boulders north
func rollNorth(grid []string) []string {
	newGrid := []string{}
	for _, line := range transpose(grid) {
		newLine := rollRow(line)
		newGrid = append(newGrid, newLine)
	}
	return transpose(newGrid)
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

	total := sumWeight(rollNorth(grid))
	fmt.Println("total", total)

}
