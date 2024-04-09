package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Hailstone struct {
	x, y, z    float64
	vx, vy, vz float64
}

func (h1 Hailstone) coefficients() (float64, float64) {
	m := h1.vy / h1.vx
	b := h1.y - (h1.x * m)
	return m, b
}

func (h1 Hailstone) WillCollide(h2 Hailstone) bool {
	// If moving apart in X, then no
	if math.Abs(h2.x-h1.x) < math.Abs(h2.x+h2.vx)-(h1.x+h1.vx) {
		return false
	}

	// If moving apart in Y, then no
	if math.Abs(h2.y-h1.y) < math.Abs(h2.y+h2.vy)-(h1.y+h1.vy) {
		return false
	}

	// if parallel in only 1 dimension, then no
	parallelInX := ((h2.x - h1.x) == (h2.x+h2.vx)-(h1.x+h1.vx)) && h1.x != h2.x
	parallelInY := ((h2.y - h1.y) == (h2.y+h2.vy)-(h1.y+h1.vy)) && h1.y != h2.y
	if parallelInX && parallelInY {
		return false
	}

	// Otherwise true
	return true
}

func (h1 Hailstone) CollidesAt(h2 Hailstone) (float64, float64) {
	// Assumes they will collide

	m1, b1 := h1.coefficients()
	m2, b2 := h2.coefficients()
	xCollision := (b2 - b1) / (m1 - m2)
	yCollision := m1*xCollision + b1
	fmt.Println("m1", m1, "b1", b1, "m2", m2, "b2", b2)
	fmt.Println("xCollision", xCollision, "yCollision", yCollision)

	return xCollision, yCollision
}

func (h1 Hailstone) CollidesInTestArea(h2 Hailstone) bool {
	if !h1.WillCollide(h2) {
		fmt.Println("No collision expected")
		return false
	}

	// baby test
	// testMin := 7.0
	// testMax := 27.0

	// full test
	testMin := 200000000000000.0
	testMax := 400000000000000.0

	xCol, yCol := h1.CollidesAt(h2)
	return xCol >= testMin && xCol <= testMax && yCol >= testMin && yCol <= testMax
}

func ProcessBlock(lines []string) int {
	// process all hailstones against each other (vs what we'd do in a test
	// program, where we're just looking at specific pairs)
	hailstones := []Hailstone{}
	for i := 0; i < len(lines); i++ {
		line := strings.Split(lines[i], "@")
		positions := strings.Split(line[0], ",")
		velocities := strings.Split(line[1], ",")
		h := Hailstone{
			x:  ParseSingle(positions[0]),
			y:  ParseSingle(positions[1]),
			z:  ParseSingle(positions[2]),
			vx: ParseSingle(velocities[0]),
			vy: ParseSingle(velocities[1]),
			vz: ParseSingle(velocities[2]),
		}
		hailstones = append(hailstones, h)
	}

	total := 0
	fmt.Println(hailstones)
	for i := 0; i < len(hailstones)-1; i++ {
		for j := i + 1; j < len(hailstones); j++ {
			res := hailstones[i].CollidesInTestArea(hailstones[j])
			fmt.Println("res", res, hailstones[i], hailstones[j])
			if res {
				total += 1
			}
		}
	}

	return total
}

func ParseSingle(s string) float64 {
	r, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return r
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)
	lines := []string{}
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			intermediateTotal := ProcessBlock(lines)
			lines = []string{}
			total += intermediateTotal
		} else {
			lines = append(lines, line)
		}
	}
	if len(lines) > 0 {
		intermediateTotal := ProcessBlock(lines)
		total += intermediateTotal
	}

	// 58560 is too high
	fmt.Println("Total", total)

}
