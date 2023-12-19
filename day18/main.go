package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Point struct {
	row, col int
	idx      int
	dir      string
}

func (p *Point) makeNext(dir string, l int) Point {
	next := Point{row: p.row, col: p.col}
	p.dir = dir
	if dir == "R" {
		next.col += l
	} else if dir == "L" {
		next.col -= l
	} else if dir == "D" {
		next.row += l
	} else if dir == "U" {
		next.row -= l
	} else {
		panic("unknown direction to move")
	}
	return next
}

func (p Point) forward(points []Point) Point {
	if p.idx+1 == len(points) {
		return points[0]
	} else {
		return points[p.idx+1]
	}
}

func (p Point) back(points []Point) Point {
	if p.idx == 0 {
		return points[len(points)-1]
	} else {
		return points[p.idx-1]
	}
}

func (p Point) RectIsInterior(points []Point) bool {
	// Check how many edge crossings left of the left edge of this rect, where p
	// is at the top-left
	numCrossings := 0
	for _, candidateP := range points {
		// Skip all horizontal edges
		if candidateP.col != candidateP.forward(points).col {
			continue
		}
		// Skip edges to the right
		if candidateP.col > p.col {
			continue
		}
		// Skip edges entirely above or below this row
		if (candidateP.row > p.row && candidateP.forward(points).row > p.row) ||
			(candidateP.row <= p.row && candidateP.forward(points).row <= p.row) {
			continue
		}
		// crossing
		// fmt.Println("found a crossing", p, candidateP, candidateP.forward(points))
		numCrossings += 1
	}
	return numCrossings%2 == 1
}

func Reindex(points []Point) []Point {
	for i := 0; i < len(points); i++ {
		points[i].idx = i
	}
	return points
}

// func (p Point) SizeOfRect(points []Point) int {

// 	if p.dir == "R" {
// 		// width determed by the forward path, going right
// 		width := p.forward(points).col - p.col
// 	} else if p.back(points).dir == "L" {
// 		// width is determined by the backwards path, going L into this point
// 		width := p.back(points).col - p.col
// 	}

// 	return 2 * width
// }

func printGrid(points []Point) {

	grid := map[[2]int]bool{}
	for _, p := range points {
		grid[[2]int{p.row, p.col}] = true
	}
	fmt.Println(grid)
	minRow := 0
	minCol := 0
	maxRow := 0
	maxCol := 0
	for p, _ := range grid {
		if p[0] > maxRow {
			maxRow = p[0]
		}
		if p[0] < minRow {
			minRow = p[0]
		}
		if p[1] > maxCol {
			maxCol = p[1]
		}
		if p[1] < minCol {
			minCol = minCol
		}
	}

	for i := minRow; i < maxRow+1; i++ {
		fmt.Printf("%03d ", i)
		for j := minCol; j < maxCol+1; j++ {
			if grid[[2]int{i, j}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//
	var curPoint Point
	var firstDir string
	points := []Point{}
	for _, line := range lines {

		// Part 1-style
		// fields := strings.Fields(line)
		// dir := fields[0]
		// numSteps, _ := strconv.Atoi(fields[1])

		// Part 2-style
		matcher := regexp.MustCompile(`..(.....)(.).`)
		fields := matcher.FindStringSubmatch(strings.Fields(line)[2])
		numSteps, _ := strconv.ParseInt(fields[1], 16, 32)
		dir := map[string]string{"0": "R", "1": "D", "2": "L", "3": "U"}[fields[2]]

		fmt.Println("input", dir, numSteps, fields)
		fmt.Println(curPoint)
		if len(points) == 0 {
			curPoint = Point{row: 0, col: 0, idx: 0, dir: dir}
			firstDir = dir
		} else {
			points[len(points)-1].dir = dir
		}
		curPoint = curPoint.makeNext(dir, int(numSteps))
		points = append(points, curPoint)
	}
	// Close the loop
	points[len(points)-1].dir = firstDir

	points = Reindex(points)
	for _, p := range points {
		fmt.Println("hmm", p)
	}

	// Add breakpoints
	rowMap := map[int]bool{}
	colMap := map[int]bool{}
	for _, p := range points {
		rowMap[p.row] = true
		colMap[p.col] = true
	}
	rows := []int{}
	for row, _ := range rowMap {
		rows = append(rows, row)
	}
	slices.Sort(rows)
	cols := []int{}
	for col, _ := range colMap {
		cols = append(cols, col)
	}
	slices.Sort(cols)
	fmt.Println("cols", cols)

	newPoints := []Point{}
	for _, p := range points {
		newPoints = append(newPoints, p)
		curPoint := p
		if p.dir == "U" {
			// Add new points at all breakpoints
			for i := len(rows) - 1; i >= 0; i-- {
				if rows[i] < p.row && rows[i] > p.forward(points).row {
					curPoint = curPoint.makeNext(p.dir, curPoint.row-rows[i])
					newPoints = append(newPoints, curPoint)
					newPoints[len(newPoints)-1].dir = p.dir
				}
			}
		}
		if p.dir == "D" {
			// Add new points at all breakpoints
			for i := 0; i < len(rows); i++ {
				if rows[i] > p.row && rows[i] < p.forward(points).row {
					curPoint = curPoint.makeNext(p.dir, rows[i]-curPoint.row)
					newPoints = append(newPoints, curPoint)
					newPoints[len(newPoints)-1].dir = p.dir
				}
			}
		}
		if p.dir == "L" {
			// Add new points at all breakpoints
			for i := len(cols) - 1; i >= 0; i-- {
				if cols[i] < p.col && cols[i] > p.forward(points).col {
					curPoint = curPoint.makeNext(p.dir, curPoint.col-cols[i])
					newPoints = append(newPoints, curPoint)
					newPoints[len(newPoints)-1].dir = p.dir
				}
			}
		}
		if p.dir == "R" {
			// Add new points at all breakpoints
			for i := 0; i < len(cols); i++ {
				if cols[i] > p.col && cols[i] < p.forward(points).col {
					curPoint = curPoint.makeNext(p.dir, cols[i]-curPoint.col)
					newPoints = append(newPoints, curPoint)
					newPoints[len(newPoints)-1].dir = p.dir
				}
			}
		}
	}
	points = Reindex(newPoints)

	// Now we look for all U/D edges on the interior. Their width is the
	// distance to the next U/D edge.
	//
	// We know there's a point there, because we've added extra points at every
	// intersection
	total := 0
	downTotals := 0
	for _, p := range points {
		if p.RectIsInterior(points) {
			closestRow := math.MaxInt64
			for _, candidateP := range points {
				if candidateP.row > p.row && candidateP.row < closestRow {
					closestRow = candidateP.row
				}
			}
			height := closestRow - p.row
			// If this is in interior space, we've already counted the top edge,
			// so adjust the height
			if p.dir != "R" {
				height -= 1
			}
			closestCol := math.MaxInt64
			for _, candidateP := range points {
				if candidateP.col > p.col && candidateP.col < closestCol && candidateP.row == p.row {
					closestCol = candidateP.col
				}
			}
			width := closestCol - p.col
			total += (height + 1) * width
			if p.dir == "R" && p.back(points).dir == "D" {
				total -= 1
			}
			fmt.Println("Area is", (height+1)*width, height, width, p, closestRow, closestCol)
		}

		// Count all D edges, since we haven't accounted for them
		if p.dir == "D" {
			edgeArea := p.forward(points).row - p.row
			if p.back(points).dir == "R" {
				edgeArea += 1
			} else if p.back(points).dir == "R" {
				edgeArea -= 1
			}
			fmt.Println("edge area", edgeArea, p, p.back(points))
			downTotals += edgeArea
		}
	}
	// printGrid(points)
	fmt.Println("total", total, downTotals, total+downTotals)

	// Sort, so we're working top-down
	// sort.Slice(points, func(i int, j int) bool {
	// 	if points[i].row != points[j].row {
	// 		return points[i].row < points[j].row
	// 	} else {
	// 		return points[i].col < points[j].col
	// 	}
	// })

	// numPoints := 0
	// for _, p := range points {
	// 	if p.RectIsInterior(points) {
	// 		numPoints += p.SizeOfRect(points)
	// 	}
	// 	// fmt.Println(p, p.forward(points), p.RectIsInterior(points))
	// }

}
