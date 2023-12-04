package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func getLineSlice(line string) ([]int, []int) {
	rest := strings.Split(strings.Split(line, ":")[1], "|")

	parseListOfNums := func(l string) []int {
		// Parse into a map to remove dupes
		m := map[int]bool{}
		spaces, _ := regexp.Compile(" +")

		for _, numStr := range spaces.Split(strings.Trim(l, " "), -1) {
			num, _ := strconv.Atoi(numStr)
			m[num] = true
		}

		// grab the keys, sort, return
		ret := []int{}
		for m, _ := range m {
			ret = append(ret, m)
		}
		sort.Ints(ret)
		return ret
	}
	return parseListOfNums(rest[0]), parseListOfNums(rest[1])
}

func getPoints(winners []int, mine []int) int {

	wMap := map[int]bool{}
	for _, w := range winners {
		wMap[w] = true
	}
	matches := 0
	for _, m := range mine {
		if wMap[m] == true {
			matches += 1
		}
	}

	return int(math.Pow(float64(2), float64(matches)-1))
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read the whole thing in, (pad the edges
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Blow up the input into a grid where the whole part number is assigned to
	// any point it occupies
	winners := make([][]int, 0)
	mine := make([][]int, 0)
	for _, line := range lines {
		w, m := getLineSlice(line)
		winners = append(winners, w)
		mine = append(mine, m)
		fmt.Println(w)
		fmt.Println(m)
		fmt.Println("--")
	}

	points := 0
	for i := 0; i < len(winners); i++ {
		points += getPoints(winners[i], mine[i])
	}

	fmt.Println("points", points)

}
