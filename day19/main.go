package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
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

func unique(input []int) []int {
	intMap := map[int]bool{}
	for _, i := range input {
		intMap[i] = true
	}
	uniqueInts := []int{}
	for i, _ := range intMap {
		uniqueInts = append(uniqueInts, i)
	}
	slices.Sort(uniqueInts)
	return uniqueInts
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

	// Figure out how many distinct ranges we have
	breakpoints := map[string][]int{} // (x m a s) -> list of breakpoints
	breakpoints["x"] = []int{1, 4001}
	breakpoints["m"] = []int{1, 4001}
	breakpoints["a"] = []int{1, 4001}
	breakpoints["s"] = []int{1, 4001}
	for _, w := range workflows {
		for _, c := range w {
			if c.op == "<" {
				breakpoints[c.field] = append(breakpoints[c.field], c.val)
			} else if c.op == ">" {
				breakpoints[c.field] = append(breakpoints[c.field], c.val+1)
			}
		}
	}
	breakpoints["x"] = unique(breakpoints["x"])
	breakpoints["m"] = unique(breakpoints["m"])
	breakpoints["a"] = unique(breakpoints["a"])
	breakpoints["s"] = unique(breakpoints["s"])
	fmt.Println(breakpoints)

	// Parse the parts
	totalSize := 0
	for ix := 1; ix < len(breakpoints["x"]); ix++ {
		x := breakpoints["x"][ix] - 1
		fmt.Println("progress", ix, len(breakpoints["x"]))
		for im := 1; im < len(breakpoints["m"]); im++ {
			m := breakpoints["m"][im] - 1
			for ia := 1; ia < len(breakpoints["a"]); ia++ {
				a := breakpoints["a"][ia] - 1
				for is := 1; is < len(breakpoints["s"]); is++ {
					s := breakpoints["s"][is] - 1
					part := Part{"x": x, "m": m, "a": a, "s": s}
					if part.process("in", workflows) {
						// fmt.Println("valid range", x, m, a, s)
						rangeSize := (breakpoints["x"][ix] - breakpoints["x"][ix-1]) *
							(breakpoints["m"][im] - breakpoints["m"][im-1]) *
							(breakpoints["a"][ia] - breakpoints["a"][ia-1]) *
							(breakpoints["s"][is] - breakpoints["s"][is-1])
						// fmt.Println("range sizes", rangeSize,
						// 	breakpoints["x"][ix]-breakpoints["x"][ix-1],
						// 	breakpoints["m"][im]-breakpoints["m"][im-1],
						// 	breakpoints["a"][ia]-breakpoints["a"][ia-1],
						// 	breakpoints["s"][is]-breakpoints["s"][is-1],
						// )
						totalSize += rangeSize
					}
				}
			}
		}
	}
	fmt.Println("total size", totalSize)

}
