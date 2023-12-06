package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type mappingData struct {
	source string
	dest   string
	ranges [][3]int // destStart, sourceStart, length
}

type mappedRange [][2]int

func getMapping(md mappingData, source int) int {
	// Do the logic, return the number
	// Fun fact: there's no need to think about the names? the input is strictly linear
	for _, r := range md.ranges {
		destStart := r[0]
		sourceStart := r[1]
		length := r[2]
		if source >= sourceStart && source < sourceStart+length {
			diff := destStart - sourceStart
			return source + diff
		}
	}
	return source
}

func getMappedRanges(md mappingData, ranges mappedRange) mappedRange {
	for _, r := range md.ranges {
		sourceStart := r[1]
		length := r[2]
		newRanges := mappedRange{}
		for _, pair := range ranges {
			start := pair[0]
			end := pair[1]
			sourceStartIsInRange := sourceStart > start && sourceStart < end-1
			sourceEndIsInRange := sourceStart+length-1 > start && sourceStart+length-1 < end-1
			if sourceStartIsInRange && sourceEndIsInRange {
				newRanges = append(newRanges, [2]int{start, sourceStart - 1})
				newRanges = append(newRanges, [2]int{sourceStart, sourceStart + length - 1})
				newRanges = append(newRanges, [2]int{sourceStart + length, end})
			} else if sourceStartIsInRange {
				newRanges = append(newRanges, [2]int{start, sourceStart - 1})
				newRanges = append(newRanges, [2]int{sourceStart, end})
			} else if sourceEndIsInRange {
				newRanges = append(newRanges, [2]int{start, sourceStart + length - 1})
				newRanges = append(newRanges, [2]int{sourceStart + length, end})
			} else {
				newRanges = append(newRanges, [2]int{start, end})
			}
		}
		ranges = newRanges
	}

	for i := 0; i < len(ranges); i++ {
		ranges[i][0] = getMapping(md, ranges[i][0])
		ranges[i][1] = getMapping(md, ranges[i][1])
	}
	return ranges
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read the whole thing in, map N map inputs
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	// Make time/distance pairs
	times := strings.Fields(strings.Split(lines[0], ":")[1])
	distances := strings.Fields(strings.Split(lines[1], ":")[1])
	timeDistPairs := [][2]int{}
	for i := 0; i < len(times); i++ {
		time, _ := strconv.Atoi(times[i])
		distance, _ := strconv.Atoi(distances[i])
		timeDistPairs = append(timeDistPairs, [2]int{time, distance})
	}

	// Figure it out!!
	totalProduct := 1.0
	for _, pair := range timeDistPairs {
		T := float64(pair[0])
		D := float64(pair[1])
		disc := math.Sqrt(T*T - 4*D)
		sol1 := math.Floor((-T+disc)/-2) + 1
		sol2 := math.Ceil((-T-disc)/-2) - 1
		totalProduct = totalProduct * (math.Abs(sol1-sol2) + 1)
		fmt.Println(T, D, disc, sol1, sol2, (math.Abs(sol1-sol2) + 1))
	}
	fmt.Println("solution", totalProduct)

}
