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
	allZ := false
	for !allZ {
		for i := 0; i < len(directions) && !allZ; i++ {
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

			allZ = true
			for _, cur := range curNodes {
				allZ = allZ && cur.name[2] == 'Z'
			}
			if steps == int(float64(steps)/1000000)*1000000 {
				fmt.Println(steps)
				fmt.Println(curNodes)
			}
		}
	}

	fmt.Println(steps)
}
