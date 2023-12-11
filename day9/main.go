package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func findParams(seq []int) int {
	nextSeq := make([]int, len(seq)-1)
	for i := 0; i < len(seq)-1; i++ {
		nextSeq[i] = seq[i+1] - seq[i]
	}
	allZeros := true
	for i := 0; i < len(nextSeq); i++ {
		allZeros = allZeros && (nextSeq[i] == 0)
	}
	finalVal := nextSeq[len(nextSeq)-1]
	if allZeros {
		// yay, we can return!
		return finalVal
	} else {
		// we must go deeper!
		return finalVal + findParams(nextSeq)
	}
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	sequences := [][]int{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		seq := []int{}
		for i := 0; i < len(line); i++ {
			num, _ := strconv.Atoi(line[i])
			seq = append(seq, num)
		}
		sequences = append(sequences, seq)
	}

	// Find the polynomial coefficients: A1 x ^ N1 + A2 x ^ N2 ... + N_len(seq)
	// - A is the value of the last line of differences
	// - N is the number of times we have to take the derivative
	total := 0
	for _, seq := range sequences {
		nextVal := findParams(seq) + seq[len(seq)-1]
		fmt.Println("Next value", nextVal)
		total += nextVal
	}
	fmt.Println("total: ", total)
}
