package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var cache = map[string]int{}

func numValidMatches(line string, nums []int, midMatch bool, depth int) int {
	// Base case
	if len(line) == 0 {
		if (len(nums) == 1 && nums[0] == 0) || len(nums) == 0 {
			// fmt.Println(strings.Repeat(" ", depth), nums, "base success")
			return 1
		} else {
			// fmt.Println(strings.Repeat(" ", depth), nums, "base failure")
			return 0
		}
	}

	// Maybe we already know the answer!
	key := fmt.Sprintf(line, nums, midMatch)
	val, ok := cache[key]
	if ok {
		return val
	}

	numSolutions := 0
	if line[0] == '?' || line[0] == '.' {
		if len(nums) == 0 {
			// Nothing left to match, but we haven't failed here. Continue on.
			// fmt.Println(strings.Repeat(" ", depth), line, numSolutions, "empty with no nums")
			numSolutions += numValidMatches(line[1:], nums, false, depth+1)
		} else if nums[0] == 0 {
			// We have a leading 0, so consume it and recurse
			// fmt.Println(strings.Repeat(" ", depth), line, nums, numSolutions, "empty with leading 0")
			numSolutions += numValidMatches(line[1:], nums[1:], false, depth+1)
		} else if !midMatch {
			// Matches must come later, so no consumptions on the nums
			// fmt.Println(strings.Repeat(" ", depth), line, nums, numSolutions, "empty with more to go")
			numSolutions += numValidMatches(line[1:], nums, false, depth+1)
		} else {
			// else, we're in the middle of a match, but this is the empty side. Do nothing.
			// fmt.Println(strings.Repeat(" ", depth), line, nums, numSolutions, "BONK empty but we're mid-match")
		}
	}
	if line[0] == '?' || line[0] == '#' {
		if len(nums) == 0 || nums[0] == 0 {
			// Invalid state, return failure
			// fmt.Println(strings.Repeat(" ", depth), line, nums, numSolutions, "BONK valid with empty nums or leading 0")
			// return numSolutions
		} else {
			// Decrement the first num and continue on
			// fmt.Println(strings.Repeat(" ", depth), line, nums, numSolutions, "valid with more to go")
			numSolutions += numValidMatches(line[1:], append([]int{nums[0] - 1}, nums[1:]...), true, depth+1)
		}
	}
	// fmt.Println(strings.Repeat(" ", depth), line, nums, "returning with", numSolutions)
	cache[key] = numSolutions
	return numSolutions
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	grid := []string{}
	groupNums := [][]int{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		grid = append(grid, line[0]+"?"+line[0]+"?"+line[0]+"?"+line[0]+"?"+line[0])
		groupNum := []int{}
		expandedNumStr := line[1] + "," + line[1] + "," + line[1] + "," + line[1] + "," + line[1]
		for _, numStr := range strings.Split(expandedNumStr, ",") {
			num, _ := strconv.Atoi(numStr)
			groupNum = append(groupNum, num)
		}
		groupNums = append(groupNums, groupNum)
	}

	// Recurse!!
	totalMatches := 0
	for i := 0; i < len(grid); i++ {
		line := grid[i]
		nums := groupNums[i]
		fmt.Println(line)
		fmt.Println(nums)
		totalMatches += numValidMatches(line, nums, false, 0)
		fmt.Println(numValidMatches(line, nums, false, 0))
	}
	fmt.Println("total", totalMatches)

}
