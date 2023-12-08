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
	nodeMatcher := regexp.MustCompile(`([A-Z]+) = \(([A-Z]+), ([A-Z]+)\)`)
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
	curNode := nodes["AAA"]
	for curNode.name != "ZZZ" {
		for _, step := range directions {
			if step == 'L' {
				curNode = nodes[curNode.left]
			} else {
				curNode = nodes[curNode.right]
			}
			steps += 1
		}
	}

	fmt.Println(steps)
}
