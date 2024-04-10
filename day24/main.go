package main

import (
	"bufio"
	"flag"
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

func Equalsish(a, b float64) bool {
	return math.Abs((a-b)/a) < 0.0000000001
}

func (h1 Hailstone) Coefficients() (float64, float64) {
	m := h1.vy / h1.vx
	b := h1.y - (h1.x * m)
	return m, b
}

func (h1 Hailstone) Distance(h2 Hailstone) float64 {
	return math.Sqrt(math.Pow(h1.x-h2.x, 2) + math.Pow(h1.y-h2.y, 2))
}

func (h1 Hailstone) TimeToCollision(xCol, yCol float64) float64 {
	// Returns how many seconds the collision will happen in
	// If the collision is in the past, then return a negative number

	xtime := (xCol - h1.x) / h1.vx
	ytime := (yCol - h1.y) / h1.vy
	if !Equalsish(xtime, ytime) {
		fmt.Println("xtime", xtime, "ytime", ytime)
		panic("Collision times don't match between x and y")
	}

	return xtime
}

func (h1 Hailstone) CollidesAt(h2 Hailstone) (float64, float64, bool) {
	// Check for parallel. If they're parallel, then no collision
	m1, _ := h1.Coefficients()
	m2, _ := h2.Coefficients()
	if Equalsish(m1, m2) {
		fmt.Println("Looks like parallel lines", m1, m2)
		return math.NaN(), math.NaN(), false
	}

	// Find the collision point, which may be in the past
	m1, b1 := h1.Coefficients()
	m2, b2 := h2.Coefficients()
	xCollision := (b2 - b1) / (m1 - m2)
	yCollision := (m1*b2 - m2*b1) / (m1 - m2)
	fmt.Println("h1", h1, "h2", h2)
	fmt.Println("m1", m1, "b1", b1, "m2", m2, "b2", b2)
	fmt.Println("xCollision", xCollision, "yCollision", yCollision)

	// If either collision point is in the past, then no collision
	time1 := h1.TimeToCollision(xCollision, yCollision)
	time2 := h2.TimeToCollision(xCollision, yCollision)
	if time1 < 0 && time2 < 0 {
		fmt.Println("Both collisions were in the past", time1, time2)
		return xCollision, yCollision, false
	}
	if time1 < 0 || time2 < 0 {
		fmt.Println("Collision was in the past", time1, time2)
		return xCollision, yCollision, false
	}

	// Collision point is in the future.
	// Check that the collision point is on both lines
	if !Equalsish(m2*xCollision+b2, yCollision) {
		fmt.Println("Checked", math.Abs(yCollision-m2*xCollision+b2), m2*xCollision+b2)
		panic("Collision point not on line 2")
	}

	return xCollision, yCollision, true
}

func (h1 Hailstone) CollidesInTestArea(h2 Hailstone, baby bool) bool {
	var testMin, testMax float64
	// baby test
	if baby {
		testMin = 7.0
		testMax = 27.0
	} else {
		// full test
		testMax = 400000000000000.0
		testMin = 200000000000000.0
	}

	xCol, yCol, doesCollide := h1.CollidesAt(h2)
	if !doesCollide {
		return false
	}
	return xCol >= testMin && xCol <= testMax && yCol >= testMin && yCol <= testMax
}

func ProcessBlock(lines []string, baby bool) int {
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
			res := hailstones[i].CollidesInTestArea(hailstones[j], baby)
			fmt.Println("res", i, j, res)
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
	babyPtr := flag.Bool("baby", false, "Run with the baby test data")
	flag.Parse()

	var f *os.File
	if *babyPtr {
		f, _ = os.Open("./baby-input.txt")
	} else {
		f, _ = os.Open("./input.txt")
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			intermediateTotal := ProcessBlock(lines, *babyPtr)
			total += intermediateTotal
			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	}
	if len(lines) > 0 {
		intermediateTotal := ProcessBlock(lines, *babyPtr)
		total += intermediateTotal
	}

	// 58560 is too high
	// 58401 is too high
	// 55832 is too high
	// 42839 is wrong
	// 44392 is wrong
	// 11098 is right!
	fmt.Println("Total", total)

}
