package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Part map[string]int

type Condition struct {
	field string
	op    string
	val   int
	next  string
}

type Workflows map[string][]Condition

func (c Condition) process(p Part) string {
	if c.op == "<" {
		if p[c.field] < c.val {
			return c.next
		}
	} else if c.op == ">" {
		if p[c.field] > c.val {
			return c.next
		}
	} else {
		return c.next
	}
	return ""
}

func (p Part) process(label string, workflows map[string][]Condition) bool {
	conditions := workflows[label]
	for _, c := range conditions {
		nextLabel := c.process(p)
		if nextLabel == "" {
			// Didn't match
			continue
		} else if nextLabel == "A" || nextLabel == "R" {
			// Special label, terminates
			return nextLabel == "A"
		} else {
			// other label, recurse
			return p.process(nextLabel, workflows)
		}
	}
	// Out of conditions and ... nothing? Seems bad
	panic("No condition matched, also no default condition")
}

func (p Part) rating() int {
	return p["x"] + p["m"] + p["a"] + p["s"]
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Parse the workflows
	workflows := Workflows{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// Done with workflows, onward!
			break
		}
		workflowMatcher := regexp.MustCompile(`([a-z]+)\{(.*),(.+)\}`)
		workflowMatcherInner := regexp.MustCompile(`([xmas])(.)([0-9]+):([A-Za-z]+)`)
		outerMatch := workflowMatcher.FindStringSubmatch(line)
		innerMatch := workflowMatcherInner.FindAllStringSubmatch(outerMatch[2], -1)
		label := outerMatch[1]
		fmt.Println(line)
		fmt.Println(outerMatch)
		fmt.Println(innerMatch)
		fmt.Println(label)
		fmt.Println("")
		conditions := []Condition{}
		for _, m := range innerMatch {
			val, _ := strconv.Atoi(m[3])
			conditions = append(conditions, Condition{m[1], m[2], val, m[4]})
		}
		conditions = append(conditions, Condition{next: outerMatch[3]})
		workflows[label] = conditions

	}

	// Parse the parts
	parts := []Part{}
	for scanner.Scan() {
		line := scanner.Text()
		partMatcher := regexp.MustCompile(`=([0-9]+)`)
		pmatch := partMatcher.FindAllStringSubmatch(line, -1)
		fmt.Println(pmatch)
		x, _ := strconv.Atoi(pmatch[0][1])
		m, _ := strconv.Atoi(pmatch[1][1])
		a, _ := strconv.Atoi(pmatch[2][1])
		s, _ := strconv.Atoi(pmatch[3][1])
		parts = append(parts, Part{"x": x, "m": m, "a": a, "s": s})

	}

	totalRating := 0
	for i, p := range parts {
		res := p.process("in", workflows)
		if res {
			fmt.Println("valid part", i, p.rating())
			totalRating += p.rating()
		} else {
			fmt.Println("invalid part", i)
		}
	}

	fmt.Println("total valid rating", totalRating)

}
