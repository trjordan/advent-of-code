package main

import (
	"bufio"
	"fmt"
	"os"
)

type Plot struct {
	row, col    int
	walkable    bool
	minDistance int
	isStart     bool
}

type PlotMap [][]*Plot

func (p Plot) updateAdjacencies(plots PlotMap) []*Plot {
	candidates := []*Plot{}
	candidates = append(candidates,
		plots[p.row-1][p.col],
		plots[p.row+1][p.col],
		plots[p.row][p.col-1],
		plots[p.row][p.col+1],
	)

	updatedPlots := []*Plot{}
	// fmt.Println("candidates", p, candidates)
	for i := 0; i < len(candidates); i++ {
		// fmt.Println("working", *candidates[i])
		// Not walkable, don't do it
		if !candidates[i].walkable {
			// fmt.Println("not walkable")
			continue
		} else if candidates[i].minDistance == 0 && !candidates[i].isStart {
			// fmt.Println("freshly marked'", p.minDistance+1)
			updatedPlots = append(updatedPlots, candidates[i])
			candidates[i].minDistance = p.minDistance + 1
		} else if candidates[i].minDistance > p.minDistance+1 {
			// fmt.Println("shorter'", p.minDistance+1)
		}
	}
	return updatedPlots
}

func isReachable(p *Plot, steps int) bool {
	if !p.isStart && p.minDistance == 0 {
		// Check for fully surrounded
		return false
	}
	if !p.walkable {
		// Nope, if not walkable at all
		return false
	}
	// Otherwise compare to number of steps and parity
	return (p.minDistance%2) == (steps%2) && p.minDistance <= steps
}

func printGrid(plots PlotMap, steps int) {
	var Reset = "\033[0m"
	var Red = "\033[31m"
	// var Green = "\033[32m"
	// var Yellow = "\033[33m"
	// var Blue = "\033[34m"
	// var Purple = "\033[35m"
	// var Cyan = "\033[36m"
	// var Gray = "\033[37m"
	// var White = "\033[97m"

	for i := 0; i < len(plots); i++ {
		for j := 0; j < len(plots[i]); j++ {
			w := plots[i][j].walkable
			if plots[i][j].isStart {
				fmt.Printf(Red + "SSSS" + Reset)
			} else if w {
				if isReachable(plots[i][j], steps) {
					fmt.Printf(Red+" %03d"+Reset, plots[i][j].minDistance)
				} else {
					fmt.Printf(" %03d", plots[i][j].minDistance)
				}
			} else {
				fmt.Print("-##-")
			}
		}
		fmt.Println()
	}
}

func countReachable(numSteps int, startPlot *Plot, plots PlotMap) int {
	// Reset the plots, we're re-using it
	for i := 0; i < len(plots); i++ {
		for j := 0; j < len(plots[0]); j++ {
			plots[i][j].minDistance = 0
			plots[i][j].isStart = false
		}
	}

	// fmt.Println("start point", startPlot)

	// Update all the distances
	queue := []*Plot{}
	queue = append(queue, startPlot)
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		toUpdate := p.updateAdjacencies(plots)
		queue = append(queue, toUpdate...)
	}

	// Manually correct the start plot
	startPlot.minDistance = 0
	startPlot.isStart = true

	// Figure out how many are reachable
	numValid := 0
	for i := 0; i < len(plots); i++ {
		for j := 0; j < len(plots[0]); j++ {
			if isReachable(plots[i][j], numSteps) {
				numValid += 1
			}
		}
	}

	return numValid
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Make the map, slightly padded
	plots := make(PlotMap, len(lines)+2)
	var startPlot Plot
	// Top/bottom margin
	plots[0] = make([]*Plot, len(lines[0])+2)
	plots[len(lines)+1] = make([]*Plot, len(lines[0])+2)
	for i := 0; i < len(lines[0])+2; i++ {
		plots[0][i] = &Plot{row: 0, col: i, walkable: false}
		plots[len(lines)+1][i] = &Plot{row: len(lines) + 1, col: i, walkable: false}
	}
	// Fill out for real, with left/right padding
	for row, line := range lines {
		plots[row+1] = make([]*Plot, len(lines)+2)
		plots[row+1][0] = &Plot{row: row + 1, col: 0, walkable: false}
		plots[row+1][len(lines[0])+1] = &Plot{row: row + 1, col: len(lines[0]) + 1, walkable: false}
		for col, chr := range line {
			plots[row+1][col+1] = &Plot{row: row + 1, col: col + 1, walkable: string(chr) == "."}
			if string(chr) == "S" {
				startPlot = *plots[row+1][col+1]
				plots[row+1][col+1].walkable = true
				plots[row+1][col+1].isStart = true
			}
		}
	}

	fmt.Println(startPlot)

	steps := 26501365
	// 487562418588924 is too low
	// 630452445283411 is naively close
	// 637087163116351 is wrong
	// 637127284466159 is wrong
	// 637128054011555 is wrong
	// 637128090024758 is wrong
	// 637128090024759 is wrong
	// 639169665638614 is too high
	// 5676619607130050 is too high

	edgeSize := len(plots) - 2
	halfSize := (edgeSize - 1) / 2 // lower -- size of 11, this is 5
	fmt.Println("plots size", edgeSize)

	corners := []*Plot{
		plots[1][1],
		plots[1][len(plots[0])-2],
		plots[len(plots)-2][1],
		plots[len(plots)-2][len(plots[0])-2],
	}
	midPoints := []*Plot{
		plots[1][halfSize+1],
		plots[halfSize+1][1],
		plots[len(plots)-2][halfSize+1],
		plots[halfSize+1][len(plots[0])-2],
	}
	fmt.Println("corners", corners)
	fmt.Println("midpoints", midPoints)

	// for steps := 7; steps < 101; steps += 2 {

	total := 0

	// 3^2 - 4 (1 * (1+1) / 2 * 4)
	// 5^2 - 12 (2 * (2+1) / 2 * 4
	// 7^2 - 24 ((7-2) * 4)

	// How many boxes from the center do we need to consider? For an 11-grid:
	// 5 steps -> 1
	// 6 - 16 steps -> 2
	// 17 - 27 steps -> 3
	// etc.
	linear := (halfSize + steps) / edgeSize

	// How many squares are fully reachable, in the same way as the starting
	// square? Meaning: if steps = row + col, you can get from anywhere on any
	// edge to any point.
	// linear 1 -> 1 full square -> 1 even square
	// linear 2 -> 5 full squares -> 1 even
	// linear 3 -> 13 full squares -> 9 even
	// linear 4 -> 25 full squares -> 9 even
	var adjustedLinear, numFullSquares, countFullSquares, totalFullSquares int
	if steps < halfSize*2 {
		// weirdo case we don't really care about, but handle it anyway
		adjustedLinear = 1
		numFullSquares = 0
		countFullSquares = countReachable(steps, &startPlot, plots)
		totalFullSquares = countFullSquares
	} else {
		adjustedLinear = 1 + (steps-2*halfSize)/(edgeSize)
		if adjustedLinear%2 == 1 {
			numFullSquares = adjustedLinear * adjustedLinear
		} else {
			numFullSquares = (adjustedLinear - 1) * (adjustedLinear - 1)

		}
		// Only need edgeSize for real input, but baby is tougher
		// Also, there's "odd" and "even" squares here
		countFullSquares = countReachable(edgeSize+steps%2+1, &startPlot, plots)
		printGrid(plots, edgeSize+steps%2+1)
		totalFullSquares = numFullSquares * countFullSquares
	}
	total += totalFullSquares
	fmt.Println("num / count of squares", numFullSquares, countFullSquares, numFullSquares*countFullSquares)

	// How many squares are fully reachable, but with the inverse odd/even
	// pattern as the starting square? Meaning: if steps = row + col, you can
	// get from anywhere on any edge to any point.
	// linear 1 -> 1 full square -> 0 odd square
	// linear 2 -> 5 full squares -> 4 odd square
	// linear 3 -> 13 full squares -> 4 odd square
	// linear 4 -> 25 full squares -> 16 odd square
	var numInverseSquares, countInverseSquares, totalInverseSquares int
	if steps < halfSize*2 {
		// weirdo case we don't really care about, but handle it anyway
		numInverseSquares = 0
		countInverseSquares = countReachable(steps, &startPlot, plots)
		totalInverseSquares = countFullSquares
	} else {
		if adjustedLinear%2 == 0 {
			numInverseSquares = adjustedLinear * adjustedLinear
		} else {
			numInverseSquares = (adjustedLinear - 1) * (adjustedLinear - 1)
		}
		// Only need edgeSize for real input, but baby is tougher
		// Also, there's "odd" and "even" squares here
		countInverseSquares = countReachable(edgeSize+steps%2, &startPlot, plots)
		printGrid(plots, edgeSize+steps%2)
		totalInverseSquares = numInverseSquares * countInverseSquares
	}
	total += totalInverseSquares
	fmt.Println(
		"num / count of inverse squares",
		numInverseSquares, countInverseSquares, numFullSquares*countInverseSquares,
	)

	// How many squares are on the edges of the diamond and identical to a
	// square that's 45 degrees from the origin?
	// 0-6 steps -> 0
	// 11-23 steps -> 1
	// 24 - 56 steps -> 3
	// etc.
	numDiagSquares := 1 + (steps-(halfSize*2+2))/(edgeSize*2)*2
	if steps < halfSize*2+2 {
		// Too small, reset
		numDiagSquares = 0
	} else if steps%edgeSize == 0 {
		// Pure intersection, all diag candidates are either 0 or counted as full
		numDiagSquares = 0
	}
	for _, corner := range corners {
		diagSteps := (steps - (edgeSize + 1)) % (edgeSize * 2)
		countDiagSquares := countReachable(diagSteps, corner, plots)
		total += numDiagSquares * countDiagSquares
		fmt.Println("num / count of diags", numDiagSquares, countDiagSquares, diagSteps, corner)
		printGrid(plots, diagSteps)
	}

	// How many squares are on the edges of the diamond and identical to a
	// square that's NOT 45 degrees from the origin?
	// 0-16 steps -> 0
	// 17-29 steps -> 2
	// 42-53 steps -> 4
	// etc.
	startingOffset := halfSize*2 + 2 + edgeSize
	numOffsetSquares := 2 + (steps-startingOffset)/(edgeSize*2)*2
	if steps < startingOffset {
		// Small enough values that we don't have any -- these are handled
		// by two-caps
		numOffsetSquares = 0
	} else if steps%edgeSize == edgeSize-1 {
		// Skimming the edges, these are counted as full
		numOffsetSquares = 0
	}
	for _, corner := range corners {
		offsetSteps := (steps - startingOffset) % (edgeSize * 2)
		countOffsetSquares := countReachable(offsetSteps, corner, plots)
		printGrid(plots, offsetSteps)
		total += numOffsetSquares * countOffsetSquares
		fmt.Println(
			"num / count of offsets",
			numOffsetSquares, countOffsetSquares, offsetSteps, corner,
		)
	}

	// Do we have 1 or 2 squares along each axis that aren't fully filled in?
	//
	// isTwoCap := (steps%edgeSize) > halfSize && ((steps+1)%edgeSize) != 0
	isTwoCap := (steps-halfSize)%edgeSize < halfSize && (steps-halfSize)%edgeSize != 0
	for _, mid := range midPoints {
		endcapSteps := (steps - (halfSize + 1)) % edgeSize
		countEndSquare := countReachable(endcapSteps, mid, plots)
		fmt.Println("endcaps steps", endcapSteps, countEndSquare, endcapSteps, mid)
		printGrid(plots, endcapSteps)
		total += countEndSquare
		if isTwoCap {
			almostEndcapSteps := endcapSteps + edgeSize
			countAlmostEndSquare := countReachable(almostEndcapSteps, mid, plots)
			// printGrid(plots, almostEndcapSteps)
			fmt.Println("2x endcaps steps", endcapSteps+edgeSize, countAlmostEndSquare, almostEndcapSteps, mid)
			total += countAlmostEndSquare
		}
	}

	// Simplified check, no barriers
	newTotal := (steps + 1) * (steps + 1)
	fmt.Println("simple check", newTotal)

	fmt.Println("final", fmt.Sprintf("%02d", steps), linear, numFullSquares, numDiagSquares, numOffsetSquares, isTwoCap, newTotal, total)
	// if newTotal != total {
	// 	panic("Failed simple check")
	// }
	// }
}
