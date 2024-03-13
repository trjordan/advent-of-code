package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Brick struct {
	x1, x2, y1, y2, z1, z2 int
	label                  string
}

type BrickNode struct {
	brick Brick
	above int
	below int
}

func parseIntUnsafe(i string) int {
	num, _ := strconv.Atoi(i)
	return num
}

func (b1 Brick) IsXYIntersecting(b2 Brick) bool {
	if b1.x2 < b2.x1 || b1.x1 > b2.x2 {
		// Fully east/west of each other
		return false
	}
	if b1.y2 < b2.y1 || b1.y1 > b2.y2 {
		// Fully north/south of each other
		return false
	}
	return true
}

func (b1 Brick) IsZIntersecting(b2 Brick) bool {
	if b1.z2 < b2.z1 || b1.z1 > b2.z2 {
		// Fully up/down of each other
		return false
	}
	return true
}

func (b1 Brick) IsRestingOn(b2 Brick) bool {
	return b1.IsXYIntersecting(b2) && b1.z1 == b2.z2+1
}

func (b Brick) IsRestingOnGround() bool {
	return b.z1 == 1
}

func (b Brick) MakeFallenBrick(newZ1 int) Brick {
	return Brick{
		x1:    b.x1,
		x2:    b.x2,
		y1:    b.y1,
		y2:    b.y2,
		z1:    newZ1,
		z2:    newZ1 + (b.z2 - b.z1),
		label: b.label,
	}
}

func (b Brick) SettleDown(below Brick) (Brick, bool) {
	// Check for badly intersecting
	if b.IsXYIntersecting(below) && b.IsZIntersecting(below) {
		//fmt.Println("uh oh", b, below, b.IsRestingOn(below))
		panic("Got into a bad state, trying to settle overlapping bricks")
	}

	// If it's already settled, no need to do anything
	if b.IsRestingOn(below) || b.IsRestingOnGround() {
		return b, false
	}

	// If XY intersecting and below is truly below, fall down
	if b.IsXYIntersecting(below) && below.z2 < b.z1 {
		//fmt.Println("Falling to intersection", b, below, b.IsRestingOn(below))
		return b.MakeFallenBrick(below.z2 + 1), true
	}

	// Otherwise go to the ground
	// This covers both where below is actually above, or they're non-intersecting
	//fmt.Println("Falling to ground", b)
	return b.MakeFallenBrick(1), true
}

func SettleAll(bricks []Brick) ([]Brick, bool) {
	somethingChanged := false
	ground := Brick{0, 0, 0, 0, 0, 0, "ground"}
	for i := 0; i < len(bricks); i++ {
		fmt.Println("Working on", bricks[i])
		bestCandidate, bestIsChange := bricks[i].SettleDown(ground)
		for j := 0; j < i; j++ {
			nextCandidate, changed := bricks[i].SettleDown(bricks[j])
			fmt.Println("Attempting settling to", bricks[j], nextCandidate)
			if !changed {
				// Found a supporting brick, we're done
				//fmt.Println("Existing support")
				bestIsChange = false
			} else if nextCandidate.z1 > bestCandidate.z1 {
				// Found a higher support than the previous
				//fmt.Println("Higher support")
				bestCandidate = nextCandidate
			}
		}
		if bestIsChange {
			//fmt.Println("Something changed, went with", bestCandidate)
			bricks[i] = bestCandidate
			somethingChanged = true
		} else {
			//fmt.Println("Did nothing", bricks[i])
		}
	}
	return bricks, somethingChanged
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	bricks := []Brick{}
	splitter := regexp.MustCompile("[,~]")
	curLabel := 'A'
	for scanner.Scan() {
		coords := splitter.Split(scanner.Text(), -1)
		bricks = append(bricks, Brick{
			x1:    parseIntUnsafe(coords[0]),
			y1:    parseIntUnsafe(coords[1]),
			z1:    parseIntUnsafe(coords[2]),
			x2:    parseIntUnsafe(coords[3]),
			y2:    parseIntUnsafe(coords[4]),
			z2:    parseIntUnsafe(coords[5]),
			label: string(curLabel),
		})
		curLabel++
	}

	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].z1 < bricks[j].z1
	})

	for _, b := range bricks {
		fmt.Println(b)
	}

	needsSettling := true
	for needsSettling {
		bricks, needsSettling = SettleAll(bricks)
		fmt.Println("NEXT ROUND ")
		for _, b := range bricks {
			fmt.Println(b)
		}
	}

	fmt.Println("Settled")
	for _, b := range bricks {
		fmt.Println(b)
	}

	aboves := make(map[Brick][]Brick)
	belows := make(map[Brick][]Brick)
	for i := 0; i < len(bricks); i++ {
		for j := 0; j < i; j++ {
			if bricks[i].IsRestingOn(bricks[j]) {
				aboves[bricks[j]] = append(aboves[bricks[j]], bricks[i])
				belows[bricks[i]] = append(belows[bricks[i]], bricks[j])
			}
		}
	}

	totalToDissolve := 0
	for i := 0; i < len(bricks); i++ {
		b := bricks[i]
		fmt.Println(b, len(aboves[b]), len(belows[b]))
		allHaveSupport := true
		for _, a := range aboves[b] {
			allHaveSupport = allHaveSupport && len(belows[a]) > 1
		}
		if allHaveSupport {
			fmt.Println("Can dissolve", b)
			totalToDissolve += 1
		}
	}
	fmt.Println("Total to dissolve", totalToDissolve)

}
