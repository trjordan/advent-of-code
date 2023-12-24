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

func isReachable(walkable bool, minDistance int, steps int) bool {
	return walkable &&
		(minDistance%2) == (steps%2) &&
		minDistance <= steps
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
				if isReachable(plots[i][j].walkable, plots[i][j].minDistance, steps) {
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

	// Figure out how many are reachable
	numValid := 0
	for i := 0; i < len(plots); i++ {
		for j := 0; j < len(plots[0]); j++ {
			if isReachable(plots[i][j].walkable, plots[i][j].minDistance, numSteps) {
				numValid += 1
			}
		}
	}

	return numValid
}

func main() {
	f, _ := os.Open("./baby-input.txt")

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

	// steps := 15
	// edgeSize := len(plots) - 2
	edgeSize := 11
	fmt.Println("plots size", edgeSize)
	for steps := 0; steps < 50; steps++ {

		// 3^2 - 4 (1 * (1+1) / 2 * 4)
		// 5^2 - 12 (2 * (2+1) / 2 * 4
		// 7^2 - 24 ((7-2) * 4)

		// How many boxes from the center do we need to consider? For an 11-grid:
		// 5 steps -> 1
		// 6 - 16 steps -> 2
		// 17 - 18 steps -> 3
		// etc.
		halfSize := (edgeSize - 1) / 2 // lower -- size of 11, this is 5
		linear := (halfSize + steps) / edgeSize

		// What is the total size of the grid we have to consider?
		// 5 steps -> 1
		// 6 - 16 steps -> 5
		// 17 - 28 steps -> 13
		fullSquare := (2*linear + 1)

		// How many squares are fully reachable? Meaning: if steps = row + col,
		// you can get from anywhere on any edge to any point.
		// fullSquare 1 -> 1 full square
		// fullSquare 3 -> 5 full squares
		// fullSquare 5 -> 13 full squares
		numFullSquares := fullSquare*fullSquare - (linear+1)*linear*2
		// countFullSquares := countReachable(edgeSize*2, &startPlot, plots) // Only need edgeSize for real input, but baby is tougher
		// fmt.Println("num / count of squares", numFullSquares, countFullSquares)

		// How many squares are on the edges of the diamond and identical to a
		// square that's 45 degrees from the origin?
		// 0-6 steps -> 0
		// 11-23 steps -> 1
		// 24 - 56 steps -> 3
		// etc.
		var numDiagSquares int
		if steps <= halfSize {
			numDiagSquares = 0
		} else {
			numDiagSquares = ((steps+halfSize+edgeSize)/2/edgeSize)*2 - 1
		}
		// countDiagSquares := countReachable(edgeSize, plots[1][1], plots)
		// fmt.Println("num / count of diags", numDiagSquares, countDiagSquares)

		// How many squares are on the edges of the diamond and identical to a
		// square that's NOT 45 degrees from the origin?
		// 0-16 steps -> 0
		// 17-29 steps -> 2
		// 30 - 52 steps -> 4
		// etc.
		numOffsetSquares := (steps - 1) / edgeSize / 2 * 2
		// countOffsetSquares := countReachable(edgeSize, &startPlot, plots)
		// fmt.Println("num / count of offsets", numOffsetSquares, countOffsetSquares)

		isTwoCap := (steps % edgeSize) > halfSize

		fmt.Println(steps, linear, numFullSquares, numDiagSquares, numOffsetSquares, isTwoCap)
	}
}
