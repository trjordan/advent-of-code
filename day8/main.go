package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type node struct {
	name  string
	left  string
	right string
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	nodes := map[string]node{}
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	directions := lines[0]
	nodeMatcher := regexp.MustCompile(`([0-9A-Z]+) = \(([0-9A-Z]+), ([0-9A-Z]+)\)`)
	for _, line := range lines[2:] {
		if strings.Trim(line, " ") == "" {
			continue
		}
		nodeMatches := nodeMatcher.FindStringSubmatch(line)
		n := node{
			name:  nodeMatches[1],
			left:  nodeMatches[2],
			right: nodeMatches[3],
		}
		nodes[n.name] = n
	}

	// Traverse!
	steps := 0
	curNodes := []node{}
	for _, n := range nodes {
		if n.name[2] == 'A' {
			curNodes = append(curNodes, n)
		}
	}
	allCycleFound := false
	origNodes := make([]node, len(curNodes))
	copy(origNodes, curNodes)
	fmt.Println("original", origNodes)

	// We're going to look for cycles individually, then do math to find the
	// cycle time of all of them at once
	cycleStarts := make([]int, len(origNodes))
	cycleLengths := make([]int, len(origNodes))
	cycleFound := make([]bool, len(origNodes))
	for !allCycleFound {
		for i := 0; i < len(directions) && !allCycleFound; i++ {

			// Progress all the nodes at once
			step := directions[i]
			newCurNodes := []node{}
			for _, curNode := range curNodes {
				if step == 'L' {
					curNode = nodes[curNode.left]
				} else {
					curNode = nodes[curNode.right]
				}
				newCurNodes = append(newCurNodes, curNode)
			}
			curNodes = newCurNodes
			steps += 1

			// Check if any of the nodes have hit a Z
			for j, cur := range curNodes {
				if cur.name[2] == 'Z' && !cycleFound[j] {
					if cycleStarts[j] == 0 {
						// First time we found the Z! Write it down
						cycleStarts[j] = steps
						fmt.Println("Found an initial Z", j, cycleLengths[j], cycleFound[j], steps)
					} else {
						// Second time! We can now know the cycle length, so
						// compute and close the trapdoop
						fmt.Println("Found a cycle", j, steps)
						cycleLengths[j] = steps - cycleStarts[j]
						cycleFound[j] = true
					}
				}
			}

			allCycleFound = true
			for _, f := range cycleFound {
				allCycleFound = allCycleFound && f
			}
			// if steps == int(float64(steps)/1000000)*1000000 {
			// 	fmt.Println(steps)
			// 	fmt.Println(curNodes)
			// }
		}
	}

	// My answer is: 16,563,603,485,021
	//fmt.Println(cycleStarts)
	fmt.Println("put the following numbers into https://www.calculatorsoup.com/calculators/math/lcm.php")
	fmt.Println(cycleLengths)
	//fmt.Println(steps)
}
