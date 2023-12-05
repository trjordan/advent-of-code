package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type mappingData struct {
	source string
	dest   string
	ranges [][3]int // destStart, sourceStart, length
}

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

	fmt.Println("hm")

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

	fmt.Println("hm")

	// Collaps the maps into seed -> location (first through last)
	locations := make([]int, len(seeds))
	copy(locations, seeds)
	for i := 0; i < len(maps); i++ {
		for j := 0; j < len(seeds); j++ {
			locations[j] = getMapping(maps[i], locations[j])
		}
	}

	fmt.Println("seeds", seeds)
	fmt.Println("locations", locations)
	fmt.Println("smallest", slices.Min(locations))
}
