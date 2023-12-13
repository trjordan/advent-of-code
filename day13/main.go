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

func findReflection(lines []string, origReflection int) int {
	// Find the initial reflection, then work our way out
	for i := 1; i < len(lines); i++ {
		currentReflection := 0 // Invalid reflection value -- implies a border at 0 / -1
		if lines[i] == lines[i-1] {
			// Reflection is at the border of i / i-1
			currentReflection = i

			for j := currentReflection; j < len(lines) && 2*i-j-1 >= 0; j++ {
				if lines[j] != lines[2*i-j-1] {
					// Something didn't match, bonk
					currentReflection = 0
					break
				}
			}
		}

		// Found something new , quit checking
		if currentReflection != 0 && currentReflection != origReflection {
			return currentReflection
		}
	}

	return 0
}

func findSmudgedReflections(grid []string, reflection int, isHorizontal bool) (int, bool) {
	// Flip all grid squares by making a new grid each time
	// Super efficient, I'm sure
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[0]); col++ {
			for _, newIsHorizontal := range []bool{false, true} {
				newGrid := make([]string, len(grid))
				copy(newGrid, grid)
				if newGrid[row][col] == '#' {
					newGrid[row] = newGrid[row][:col] + "." + newGrid[row][col+1:]
				} else {
					newGrid[row] = newGrid[row][:col] + "#" + newGrid[row][col+1:]
				}

				// Check both directions in this new grid
				// Check if we have a new reflection
				var newReflection int
				if !newIsHorizontal {
					newGrid = transpose(newGrid)
				}
				if isHorizontal == newIsHorizontal {
					// Make sure we avoid returing the original one
					newReflection = findReflection(newGrid, reflection)
				} else {
					newReflection = findReflection(newGrid, -1)

				}

				// Only update if this isn't the existing reflection
				isSameReflection := newReflection == reflection && isHorizontal == newIsHorizontal
				if newReflection != 0 && !isSameReflection {
					// fmt.Println("New reflection!", newReflection, row, col, newIsHorizontal)
					// fmt.Println("original grid")
					// printGrid(grid)
					// fmt.Println("new grid")
					// printGrid(newGrid)
					return newReflection, newIsHorizontal
				}
			}
		}
	}
	fmt.Println("Found", reflection, isHorizontal)
	printGrid(grid)
	panic("Could not find a new reflection")
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	grids := [][]string{}
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			grids = append(grids, lines)
			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	}
	grids = append(grids, lines)

	// Summarize
	summary := 0
	for _, grid := range grids {
		// fmt.Println("Checking vertical", i)
		reflection := findReflection(transpose(grid), -1)
		isHorizontal := false
		if reflection == 0 {
			reflection = findReflection(grid, -1)
			// fmt.Println("Checking horizontal", i)
			isHorizontal = true
		}

		reflection, isHorizontal = findSmudgedReflections(grid, reflection, isHorizontal)

		// Add up to the thing
		if isHorizontal {
			reflection = reflection * 100
		} else {
			reflection = reflection
		}
		summary += reflection
	}

	fmt.Println("total", summary)
}
