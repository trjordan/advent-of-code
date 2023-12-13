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

func findReflection(lines []string) int {
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

		// Found something, quit checking
		if currentReflection != 0 {
			fmt.Println("found a reflection", currentReflection)
			return currentReflection
		}
	}

	return 0
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
	for i, grid := range grids {
		fmt.Println("Checking horizontal", i)
		reflection := findReflection(transpose(grid))
		if reflection == 0 {
			fmt.Println("Checking vertical", i)
			reflection = findReflection(grid) * 100
		}
		summary += reflection
	}

	fmt.Println("total", summary)
}
