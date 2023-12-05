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

func getElfMap(md mappingData) map[int]int {
	// Do the logic, return a single transformation
	// Fun fact: there's no need to think about the names? the input is strictly linear

	elfMap := map[int]int{}
	fmt.Println(md.source, md.dest)
	for _, r := range md.ranges {
		fmt.Println(r)
		destStart := r[0]
		sourceStart := r[1]
		length := r[2]

		diff := destStart - sourceStart
		for i := sourceStart; i < sourceStart+length; i++ {
			elfMap[i] = i + diff
		}
	}

	return elfMap
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
	mapLines := []mappingData{}
	var mapLinesSubset mappingData
	for i := 2; i < len(lines); i++ {
		line := lines[i]
		if strings.HasSuffix(line, "map:") {
			// Start a new map
			names := strings.Split(strings.Split(line, " ")[0], "-")
			mapLinesSubset = mappingData{
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
			mapLinesSubset.ranges = append(mapLinesSubset.ranges, ranges)
		} else {
			// We're done, append it to the list
			mapLines = append(mapLines, mapLinesSubset)
		}
	}
	// No newline at the end ef the file
	mapLines = append(mapLines, mapLinesSubset)

	fmt.Println("hm")

	// Get a list of maps
	elfMaps := []map[int]int{}
	for _, md := range mapLines {
		elfMaps = append(elfMaps, getElfMap(md))
	}

	fmt.Println("hm")

	// Collaps the maps into seed -> location (first through last)
	locations := make([]int, len(seeds))
	copy(locations, seeds)
	for i := 0; i < len(elfMaps); i++ {
		for j := 0; j < len(seeds); j++ {
			if elfMaps[i][locations[j]] != 0 {
				locations[j] = elfMaps[i][locations[j]]
			}
		}
	}

	fmt.Println("seeds", seeds)
	fmt.Println("locations", locations)
	fmt.Println("smallest", slices.Min(locations))
}
