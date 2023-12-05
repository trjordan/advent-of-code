package main

import (
	"bufio"
	"fmt"
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
	// Parse the seeds
	seedStrs := strings.Fields(strings.Split(lines[0], ":")[1])
	seeds := []int{}
	for _, seedStr := range seedStrs {
		seed, _ := strconv.Atoi(seedStr)
		seeds = append(seeds, seed)
	}

	// Parse the mapping data
	maps := []mappingData{}
	var mapsSubset mappingData
	for i := 2; i < len(lines); i++ {
		line := lines[i]
		if strings.HasSuffix(line, "map:") {
			// Start a new map
			names := strings.Split(strings.Split(line, " ")[0], "-")
			mapsSubset = mappingData{
				source: names[0],
				dest:   names[2],
				ranges: [][3]int{},
			}
		} else if line != "" {
			// Add to the existing range
			ranges := [3]int{}
			for i, n := range strings.Fields(line) {
				ranges[i], _ = strconv.Atoi(n)
			}
			mapsSubset.ranges = append(mapsSubset.ranges, ranges)
		} else {
			// We're done, append it to the list
			maps = append(maps, mapsSubset)
		}
	}
	// No newline at the end ef the file
	maps = append(maps, mapsSubset)

	// Create ranges of interest
	seedRanges := make(mappedRange, 0)
	for i := 0; i < len(seeds); i += 2 {
		seedRanges = append(seedRanges, [2]int{seeds[i], seeds[i] + seeds[i+1] - 1})
	}
	for i := 0; i < len(maps); i++ {
		seedRanges = getMappedRanges(maps[i], seedRanges)
	}

	// Find the lowest and interate
	lowestLocation := int(1e12)
	for i := 0; i < len(seedRanges); i++ {
		if seedRanges[i][0] < lowestLocation {
			lowestLocation = seedRanges[i][0]
		}
	}
	fmt.Println("lowest", lowestLocation)

}
