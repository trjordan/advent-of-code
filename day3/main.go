package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func updateLineMatches(line string, validSpaces []bool) []bool {
	validMatches := regexp.MustCompile("[^.0-9]")

	prevMatches := validMatches.FindAllStringIndex(line, -1)
	for _, match := range prevMatches {
		i := match[0]
		if i > 0 {
			validSpaces[i-1] = true
		}
		validSpaces[i] = true
		if i < len(line)-1 {
			validSpaces[i+1] = true
		}
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
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	blankline := strings.Repeat(".", len(lines[0]))
	parts := make([][]int, 0)
	for i, line := range lines {
		lineParts := make([]int, 0)
		if i == 0 {
			lineParts = getParts(blankline, line, lines[i+1])
		} else if i == len(lines)-1 {
			lineParts = getParts(lines[i-1], line, blankline)

		} else {
			lineParts = getParts(lines[i-1], line, lines[i+1])
		}
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
