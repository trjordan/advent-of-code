package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type lens struct {
	label       string
	focalLength int
}

func getHash(label string) int {
	cur := 0
	// Walk bytes, not runes
	for i := 0; i < len(label); i++ {
		cur = update(cur, int(label[i]))
	}
	return cur
}

func update(cur int, chr int) int {
	return ((cur + chr) * 17) % 256
}

func hasLabel(lenses []lens, label string) bool {
	bucketHasLabel := false
	for _, lens := range lenses {
		bucketHasLabel = bucketHasLabel || lens.label == label
	}
	return bucketHasLabel
}

func addLens(lenses []lens, label string, focalLength int) []lens {
	newLenses := []lens{}
	newLens := lens{label: label, focalLength: focalLength}
	lensWasFound := false
	for _, lens := range lenses {
		if lens.label == label {
			newLenses = append(newLenses, newLens)
			lensWasFound = true
		} else {
			newLenses = append(newLenses, lens)
		}
	}
	if !lensWasFound {
		newLenses = append(newLenses, newLens)
	}
	return newLenses
}

func removeLens(lenses []lens, label string) []lens {
	newLenses := []lens{}
	for _, lens := range lenses {
		if lens.label != label {
			newLenses = append(newLenses, lens)
		}
	}
	return newLenses
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Just one line today!
	steps := []string{}
	for scanner.Scan() {
		steps = strings.Split(scanner.Text(), ",")
	}

	// int -> slice of lenses
	boxes := map[int][]lens{}
	for i := 0; i < 256; i++ {
		boxes[i] = []lens{}
	}

	labelMatcher := regexp.MustCompile("(.+)([=-])(.*)")
	for _, step := range steps {
		// Find the box to go in
		stepParts := labelMatcher.FindStringSubmatch(step)
		fmt.Println("Working on ", step, stepParts)
		label := stepParts[1]
		hashBucket := getHash(label)
		if stepParts[2] == "=" {
			// put a new item in the list
			focalLength, _ := strconv.Atoi(stepParts[3])
			// Check if this is in the list
			if hasLabel(boxes[hashBucket], label) {
				fmt.Println("Updating label", hashBucket, label, focalLength)
				boxes[hashBucket] = addLens(boxes[hashBucket], label, focalLength)
			} else {
				fmt.Println("Adding label", hashBucket, label, focalLength)
				boxes[hashBucket] = addLens(boxes[hashBucket], label, focalLength)
			}
		} else if stepParts[2] == "-" {
			// remove this label
			if hasLabel(boxes[hashBucket], label) {
				fmt.Println("Removing label", hashBucket, label)
				boxes[hashBucket] = removeLens(boxes[hashBucket], label)
			} else {
				fmt.Println("Nothing to remove", hashBucket, label)
			}
		}
	}

	totalPower := 0
	for boxNum, lenses := range boxes {
		for lensNum, lens := range lenses {
			lensPower := (1 + boxNum) * (1 + lensNum) * lens.focalLength
			totalPower += lensPower
			fmt.Println("Power of ", lens, lensPower)
		}
	}

	fmt.Println(boxes)

	fmt.Println("total power", totalPower)
}
