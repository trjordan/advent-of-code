package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type Segment struct {
	start, end Point
	points     []Point
	length     int
}

func NextPoints(strMap []string, prev Point, cur Point) []Point {
	nextPoints := []Point{}

	// Look right
	if cur.x > 0 {
		if !(cur.x+1 == prev.x && cur.y == prev.y) &&
			(strMap[cur.y][cur.x+1] == '.' ||
				strMap[cur.y][cur.x+1] == '>') {
			nextPoints = append(nextPoints, Point{cur.x + 1, cur.y})
		}
	}

	// Look left
	if cur.x < len(strMap[0])-1 {
		if !(cur.x-1 == prev.x && cur.y == prev.y) &&
			(strMap[cur.y][cur.x-1] == '.' ||
				strMap[cur.y][cur.x-1] == '<') {
			nextPoints = append(nextPoints, Point{cur.x - 1, cur.y})
		}
	}

	// Look up
	if cur.y > 0 {
		if !(cur.x == prev.x && cur.y-1 == prev.y) &&
			(strMap[cur.y-1][cur.x] == '.' ||
				strMap[cur.y-1][cur.x] == '^') {
			nextPoints = append(nextPoints, Point{cur.x, cur.y - 1})
		}
	}

	// Look down
	if cur.y < len(strMap)-1 {
		if !(cur.x == prev.x && cur.y+1 == prev.y) &&
			(strMap[cur.y+1][cur.x] == '.' ||
				strMap[cur.y+1][cur.x] == 'v') {
			nextPoints = append(nextPoints, Point{cur.x, cur.y + 1})
		}
	}

	return nextPoints
}

func CreateSegments(strMap []string) []Segment {
	// Manual start
	nextStarts := []Point{{1, 0}, {1, 1}}
	prev := nextStarts[0]
	cur := Point{1, 1}
	segments := []Segment{}
	curPoints := []Point{}

	// Don't get stuck looping, keep track of the intersections we've seen, and
	// don't add them to nextStarts
	seen := map[Point]bool{}

	// Go explore!
	for len(nextStarts) > 0 {
		curPoints = append(curPoints, prev)
		nexts := NextPoints(strMap, prev, cur)
		if len(nexts) != 1 {
			// wrap it up
			curPoints = append(curPoints, cur)
			segments = append(segments, Segment{
				start:  nextStarts[0],
				end:    cur,
				points: curPoints,
				length: len(curPoints),
			})
			curPoints = []Point{}
			nextStarts = nextStarts[2:]

			if len(nexts) > 1 {
				for i := 0; i < len(nexts); i++ {
					if !seen[nexts[i]] {
						nextStarts = append(nextStarts, cur, nexts[i])
						seen[nexts[i]] = true
					}
				}
			}
			if len(nextStarts) > 0 {
				prev = nextStarts[0]
				cur = nextStarts[1]
			}
		} else {
			// keep chugging
			prev = cur
			cur = nexts[0]
		}
	}

	return segments
}

func FindPreviousSegments(p Point, segments []Segment) []Segment {
	segmentMap := map[Point][]Segment{}
	for i := 0; i < len(segments); i++ {
		segmentMap[segments[i].end] = append(segmentMap[segments[i].end], segments[i])
	}

	return segmentMap[p]
}

func IsStartingSegment(s Segment) bool {
	return s.start.x == 1 && s.start.y == 0
}

func FindLongestPathLen(start Point, lenSoFar int, segments []Segment, strMap []string) int {
	fmt.Println(start, lenSoFar)
	nextSegments := FindPreviousSegments(start, segments)

	longest := lenSoFar
	for i := 0; i < len(nextSegments); i++ {
		// fmt.Println("Gonna do", nextSegments[i], longest)
		// PrintMap(strMap, []Segment{nextSegments[i]})
		nextLenSoFar := FindLongestPathLen(nextSegments[i].start, lenSoFar+nextSegments[i].length-1, segments, strMap)
		if nextLenSoFar > longest {
			longest = nextLenSoFar
		}
		if IsStartingSegment(nextSegments[i]) {
			fmt.Println("Found a full path with length ", nextLenSoFar)
		}
	}

	return longest
}

func PrintMap(strMap []string, segments []Segment) {
	// Invert the segments into a map for easy lookup
	walkedPoints := map[Point]bool{}
	for _, s := range segments {
		for _, p := range s.points {
			walkedPoints[p] = true
		}
	}
	for i := 0; i < len(strMap); i++ {
		for j := 0; j < len(strMap[i]); j++ {
			if walkedPoints[Point{j, i}] {
				fmt.Print("O")
			} else {
				fmt.Printf("%c", strMap[i][j])
			}
		}
		print("\n")
	}
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	strMap := []string{}
	for scanner.Scan() {
		strMap = append(strMap, scanner.Text())
	}

	segments := CreateSegments(strMap)
	// for i := 0; i < len(segments); i++ {
	// 	fmt.Println(segments[i])
	// 	PrintMap(strMap, []Segment{segments[i]})
	// }

	// Find all the segment that end, walk backwards to find the longest
	longest := FindLongestPathLen(Point{len(strMap[0]) - 2, len(strMap) - 1}, 0, segments, strMap)
	fmt.Println("longest", longest)
	// PrintMap(strMap, segments)
}
