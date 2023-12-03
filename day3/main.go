package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// func collapsePartNums(line []string) []string {
// 	partMatches := regexp.MustCompile("[0-9]+")

// }

func updateLineMatches(line string, validSpaces []bool) []bool {
	validMatches := regexp.MustCompile("[^.0-9]")

	prevMatches := validMatches.FindAllStringIndex(line, -1)
	for _, match := range prevMatches {
		// Padding the input means we know there's no matches on the edge
		i := match[0]
		validSpaces[i-1] = true
		validSpaces[i] = true
		validSpaces[i+1] = true
	}

	return validSpaces
}

func getParts(prev string, line string, next string) []int {

	partNums := []int{}

	validSpaces := make([]bool, len(line))
	validSpaces = updateLineMatches(prev, validSpaces)
	validSpaces = updateLineMatches(line, validSpaces)
	validSpaces = updateLineMatches(next, validSpaces)

	partMatches := regexp.MustCompile("[0-9]+")

	validParts := partMatches.FindAllStringIndex(line, -1)
	for _, pair := range validParts {
		num := line[pair[0]:pair[1]]
		numInt, _ := strconv.Atoi(num)
		valid := false
		for i := pair[0]; i < pair[1]; i++ {
			valid = valid || validSpaces[i]
		}
		if valid {
			partNums = append(partNums, numInt)
		}
	}

	fmt.Println(prev)
	fmt.Println(line)
	fmt.Println(next)
	fmt.Println(validSpaces)
	fmt.Println(validParts)
	fmt.Println(partNums)
	fmt.Println("--")

	return partNums
}

func main() {
	f, _ := os.Open("./baby-input.txt")

	scanner := bufio.NewScanner(f)

	// Read the whole thing in, pad the edges
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, "."+scanner.Text()+".")
	}
	blankline := strings.Repeat(".", len(lines[0]))
	lines = append([]string{blankline}, lines...)
	lines = append(lines, blankline)

	// Go find the part numbers
	parts := make([][]int, 0)
	for i := 1; i < len(lines)-1; i++ {
		line := lines[i]
		lineParts := getParts(lines[i-1], line, lines[i+1])
		parts = append(parts, lineParts)
		fmt.Println(lineParts)
	}

	// Flatten and sum
	s := 0
	for _, ps := range parts {
		for _, p := range ps {
			s += p
		}
	}
	fmt.Println(s)
}
