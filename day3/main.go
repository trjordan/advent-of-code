package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getLineSlice(line string) []string {
	lineSlice := make([]string, len(line))

	// Parts get filled with the full part number
	partMatches := regexp.MustCompile("[0-9]+")
	partPositions := partMatches.FindAllStringIndex(line, -1)
	for _, pair := range partPositions {
		num := line[pair[0]:pair[1]]
		for i := pair[0]; i < pair[1]; i++ {
			lineSlice[i] = num
		}
	}

	// Gears get inserted literally (lol regexp overkill)
	gearMatches := regexp.MustCompile("\\*")
	gearPositions := gearMatches.FindAllStringIndex(line, -1)
	for _, pair := range gearPositions {
		for i := pair[0]; i < pair[1]; i++ {
			lineSlice[i] = "*"
		}
	}

	return lineSlice
}

func getGearRatioProducts(prev []string, line []string, next []string) []int {

	ratios := []int{}

	// Find all unique part numbers around the gear
	// (Bug: what if there are two of the same part? )
	for i := 1; i < len(line)-1; i++ {
		target := line[i]
		if target == "*" {
			partNums := []string{
				prev[i-1],
				prev[i],
				prev[i+1],
				"",
				line[i-1],
				line[i+1],
				"",
				next[i-1],
				next[i],
				next[i+1],
			}

			prev := ""
			current := ""
			parts := []int{}
			for _, partNum := range partNums {
				current = partNum
				if prev != current && current != "" {
					currentInt, _ := strconv.Atoi(current)
					parts = append(parts, currentInt)
				}
				prev = current
			}
			if len(parts) == 2 {
				ratio := 1
				for _, p := range parts {
					ratio = ratio * p
				}
				ratios = append(ratios, ratio)
			}
		}
	}

	return ratios
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read the whole thing in, pad the edges
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, "."+scanner.Text()+".")
	}
	blankline := strings.Repeat(".", len(lines[0]))
	lines = append([]string{blankline}, lines...)
	lines = append(lines, blankline)

	// Blow up the input into a grid where the whole part number is assigned to
	// any point it occupies
	parts := make([][]string, 0)
	for _, line := range lines {
		parts = append(parts, getLineSlice(line))
	}

	// Go find the ratios
	sumOfRatios := 0
	for i := 1; i < len(parts)-1; i++ {
		part := parts[i]
		ratios := getGearRatioProducts(parts[i-1], part, parts[i+1])
		for _, ratio := range ratios {
			sumOfRatios += ratio
		}
	}

	fmt.Println("final", sumOfRatios)

}
